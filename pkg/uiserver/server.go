package uiserver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gqgs/go-zeronet/cmd/site"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
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

	r.Route("/{site:1[0-9A-Za-z]{31,33}}", func(r chi.Router) {
		r.Get("/", s.SiteHandler)
		r.Get("/{file}", s.SiteFileHandler)
	})
	return s
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

func (s *server) SiteHandler(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	if err := s.siteManager.ReadFile(site, "index.html", w); err != nil {
		s.log.WithField("site", site).Warn(err)
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (s *server) SiteFileHandler(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	file := chi.URLParam(r, "file")
	if err := s.siteManager.ReadFile(site, file, w); err != nil {
		s.log.WithField("site", site).WithField("file", file).Warn(err)
		http.Error(w, "not found", http.StatusNotFound)
	}
}
