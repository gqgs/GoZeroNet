package uiserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/websocket"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/uimedia"
)

type server struct {
	srv         *http.Server
	log         log.Logger
	addr        string
	siteManager site.SiteManager
}

func NewServer(addr string, siteManager site.SiteManager) *server {
	r := chi.NewRouter()

	s := &server{
		addr: addr,
		srv: &http.Server{
			Addr:    addr,
			Handler: r,
		},
		log:         log.New("uiserver"),
		siteManager: siteManager,
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
	})
	return s
}

func (s *server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		s.log.Error(err)
		return
	}

	go newUIWebsocket(conn, s.siteManager).Serve()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D")
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
	innerPath = strings.TrimSuffix(innerPath, "/")

	if innerPath == "" {
		if err := s.siteManager.RenderIndex(site, "index.html", w); err != nil {
			s.log.WithField("site", site).Warn(err)
			http.Error(w, "not found", http.StatusNotFound)
		}
		return
	}

	if strings.HasSuffix(innerPath, "all.css") {
		w.Header().Add("Content-Type", "text/css")
	} else if strings.HasSuffix(innerPath, "all.js") {
		w.Header().Add("Content-Type", "application/javascript")
	}

	if err := s.siteManager.ReadFile(site, innerPath, w); err != nil {
		s.log.WithField("site", site).WithField("innerPath", innerPath).Warn(err)
		http.Error(w, "not found", http.StatusNotFound)
	}
}
