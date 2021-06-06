package ui

import (
	"context"
	"net/http"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
)

type server struct {
	srv *http.Server
	log log.Logger
}

func NewServer() *server {
	mux := http.NewServeMux()
	s := &server{
		srv: &http.Server{
			Addr:    config.UIServer.Addr(),
			Handler: mux,
		},
		log: log.New("uiserver"),
	}
	mux.HandleFunc("/ping", s.Ping)
	return s
}

func (s *server) Shutdown(ctx context.Context) error {
	if s == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *server) Listen() {
	s.log.Infof("listening at %s", config.UIServer.Addr())
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Fatal(err)
	}
}

func (s *server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!\n"))
}
