package ui

import (
	"context"
	"log"
	"net/http"

	"github.com/gqgs/go-zeronet/pkg/config"
)

type Server struct {
	srv *http.Server
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *Server) Listen() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong!\n"))
	})

	srv := http.Server{
		Addr:    config.UIServer.Addr(),
		Handler: mux,
	}
	s.srv = &srv

	println("ui server listening...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
