package uiserver

import (
	"archive/zip"
	"context"
	"errors"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/lib/websocket"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/uimedia"
	"github.com/gqgs/go-zeronet/pkg/uiwebsocket"
	"github.com/gqgs/go-zeronet/pkg/user"
)

type server struct {
	ctx           context.Context
	srv           *http.Server
	log           log.Logger
	addr          string
	siteManager   site.Manager
	fileServer    fileserver.Server
	pubsubManager pubsub.Manager
	userManager   user.Manager
}

func NewServer(ctx context.Context, addr string, siteManager site.Manager, fileServer fileserver.Server,
	pubsubManager pubsub.Manager, userManager user.Manager) (*server, error) {
	r := chi.NewRouter()

	host, portString, _ := net.SplitHostPort(addr)
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}

	config.UIServerHost = host
	config.UIServerPort = port

	s := &server{
		ctx:  ctx,
		addr: addr,
		srv: &http.Server{
			Addr:    addr,
			Handler: r,
		},
		log:           log.New("uiserver"),
		siteManager:   siteManager,
		fileServer:    fileServer,
		pubsubManager: pubsubManager,
		userManager:   userManager,
	}

	uimediaHandler := http.FileServer(http.FS(uimedia.FS)).ServeHTTP
	r.Get("/", indexHandler)
	r.Get("/uimedia/{file}", uimediaHandler)
	r.Get("/uimedia/img/{file}", uimediaHandler)
	r.Get("/uimedia/lib/{file}", uimediaHandler)
	r.Route("/{site:1[0-9A-Za-z]{31,33}}", func(r chi.Router) {
		r.Get("/*", s.siteHandler)
	})
	r.Route("/ZeroNet-Internal", func(r chi.Router) {
		r.Get("/Websocket", s.websocketHandler)
		r.Group(func(r chi.Router) {
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:   []string{"null"},
				AllowCredentials: true,
			}))
			r.Post("/BigfileUpload", s.bigfileUpload)
		})
		r.Get("/BigfileUploadWebsocket", s.bigfileUploadWebsocket)
	})
	return s, nil
}

func (s *server) bigfileUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(config.MultipartFormMaxMemory); err != nil {
		s.log.Error(err)
		return
	}

	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		return
	}

	uploadNonce := r.Form.Get("upload_nonce")
	upload, err := s.siteManager.FindPendingUpload(uploadNonce)
	if err != nil {
		s.log.Error(err)
		return
	}

	filePath := path.Join(config.DataDir, upload.Site, safe.CleanPath(upload.InnerPath))
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		s.log.Error(err)
		return
	}
	file, err := os.Create(filePath)
	if err != nil {
		s.log.Error(err)
		return
	}
	defer file.Close()

	for _, formFiles := range r.MultipartForm.File {
		if len(formFiles) == 0 {
			return
		}
		formFile := formFiles[0]
		if formFile.Size != int64(upload.Size) {
			s.log.Error("incorrect size")
			return
		}
		f, err := formFile.Open()
		if err != nil {
			s.log.Error(err)
			return
		}
		defer f.Close()
		if _, err := io.Copy(file, f); err != nil {
			s.log.Error(err)
		}

		return
	}
}

func (s *server) bigfileUploadWebsocket(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.log.Error(err)
		return
	}
	uploadNonce := r.Form.Get("upload_nonce")
	upload, err := s.siteManager.FindPendingUpload(uploadNonce)
	if err != nil {
		s.log.Error(err)
		return
	}

	filePath := path.Join(config.DataDir, upload.Site, safe.CleanPath(upload.InnerPath))
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		s.log.Error(err)
		return
	}
	file, err := os.Create(filePath)
	if err != nil {
		s.log.Error(err)
		return
	}
	defer file.Close()

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		s.log.Error(err)
		return
	}

	var read int
	for {
		if err := conn.Write([]byte("poll")); err != nil {
			s.log.Error(err)
			break
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			s.log.Error(err)
			break
		}

		read += len(message)

		if _, err := file.Write(message); err != nil {
			s.log.Error(err)
			break
		}

		if read >= upload.Size {
			conn.Write([]byte("done"))
			break
		}
	}
}

func (s *server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.log.Error(err)
		return
	}
	wrapperKey := r.Form.Get("wrapper_key")
	site := s.siteManager.SiteByWrapperKey(wrapperKey)
	if site == nil {
		http.Error(w, "site not found", http.StatusNotFound)
		return
	}

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		s.log.Error(err)
		return
	}

	go uiwebsocket.NewUIWebsocket(s.ctx, conn, s.siteManager, s.fileServer, site, s.pubsubManager).Serve()
	go site.Announce(site.AnnounceTrackers)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", config.HomeSite)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (s *server) Shutdown(ctx context.Context) error {
	if s == nil || s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *server) Listen(ctx context.Context) {
	s.log.Infof("listening at http://%s", s.addr)
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Fatal(err)
	}
}

func (s *server) siteHandler(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	innerPath := strings.TrimPrefix(r.URL.Path, "/"+site)
	innerPath = strings.Trim(innerPath, "/")

	if s.siteManager.Site(site) == nil {
		newSite, err := s.siteManager.NewSite(site)
		if err != nil {
			s.log.Error(err)
			return
		}
		newSite.Loading(true)

		msgCh := s.pubsubManager.Register("site_downloader", config.DefaultChannelSize)
		defer s.pubsubManager.Unregister(msgCh)

		go newSite.Announce()
		go func() {
			defer newSite.Loading(false)

			if err := newSite.Download(time.Now().AddDate(0, 0, -7)); err != nil {
				s.log.Error(err)
				return
			}

			if err := newSite.OpenDB(); err != nil {
				s.log.Error(err)
				return
			}
			if err := newSite.RebuildDB(); err != nil {
				s.log.Error(err)
				return
			}
		}()

		s.waitForContentDownload(site, msgCh)
	}

	if innerPath == "" {
		if err := s.siteManager.RenderIndex(site, "index.html", w); err != nil {
			s.log.WithField("site", site).Warn(err)
			http.Error(w, "not found", http.StatusNotFound)
		}
		return
	}

	switch {
	case strings.HasSuffix(innerPath, ".svg"):
		w.Header().Add("Content-Type", "image/svg+xml")
	case strings.HasSuffix(innerPath, ".css"):
		w.Header().Add("Content-Type", "text/css")
	case strings.HasSuffix(innerPath, ".js"):
		w.Header().Add("Content-Type", "application/javascript")
	}

	if i := strings.Index(innerPath, ".zip/"); i > 0 {
		i += len(".zip/")
		zipPath, filename := innerPath[:i], innerPath[i:]
		s.handleZip(w, r, site, zipPath, filename)
		return
	}

	if err := s.siteManager.ReadFile(site, innerPath, w, r); err != nil {
		s.log.WithField("site", site).WithField("innerPath", innerPath).Warn(err)
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *server) waitForContentDownload(site string, msgCh <-chan pubsub.Message) {
	for {
		select {
		case msg := <-msgCh:
			if msg.Site() != site {
				continue
			}
			if payload, ok := msg.Event().(*event.ContentInfo); ok {
				if payload.InnerPath == "content.json" {
					s.log.Info("downloaded content.json")
					return
				}
			}
		case <-time.After(5 * time.Minute):
			return
		}
	}
}

func (s *server) handleZip(w http.ResponseWriter, r *http.Request, site, zipPath, filename string) {
	zipHandler := &zipFS{w, r, filename}
	if err := s.siteManager.ReadFile(site, zipPath, zipHandler, nil); err != nil {
		s.log.WithField("site", site).WithField("zipPath", zipPath).Warn(err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
}

type zipFS struct {
	w        http.ResponseWriter
	r        *http.Request
	filename string
}

type fileReaderAt interface {
	io.ReaderAt
	fs.File
}

// Implements the io.ReaderFrom interface
func (z *zipFS) ReadFrom(r io.Reader) (n int64, err error) {
	if file, ok := r.(fileReaderAt); ok {
		stat, err := file.Stat()
		if err != nil {
			return 0, err
		}

		zipReader, err := zip.NewReader(file, stat.Size())
		if err != nil {
			return 0, err
		}

		z.r.URL.Path = "/" + z.filename
		http.FileServer(http.FS(zipReader)).ServeHTTP(z.w, z.r)
		return 0, nil
	}
	return 0, errors.New("unexpected reader interface")
}

// Implements the io.Writer interface
func (z *zipFS) Write(data []byte) (n int, err error) {
	return 0, errors.New("unexpected write call")
}
