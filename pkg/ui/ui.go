package ui

import (
	"context"
	"net/http"

	"github.com/gqgs/go-zeronet/pkg/lib/log"
)

type server struct {
	srv  *http.Server
	log  log.Logger
	addr string
}

func NewServer(addr string) *server {
	mux := http.NewServeMux()
	s := &server{
		addr: addr,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		log: log.New("uiserver"),
	}
	mux.HandleFunc("/ping", s.Ping)
	return s
}

func (s *server) Shutdown(ctx context.Context) error {
	if s == nil || s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *server) Listen(ctx context.Context) {
	s.log.Infof("listening at %s", s.addr)
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Fatal(err)
	}
}

func (s *server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!\n"))
}
