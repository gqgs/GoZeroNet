package file

import (
	"context"
	"log"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

// Server implements the protocol specified at:
// https://zeronet.io/docs/help_zeronet/network_protocol/
//
// Every message is encoded with MessagePack
// Every request has 3 parameters: `cmd`, `req_id` and `params`
type Server struct {
	srv *http.Server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

type requestPayload struct {
	CMD    string
	ReqID  int
	Params string
}

type pingResponsePayload struct {
	CMD  string
	To   int
	Body string
}

func pingHandler(w http.ResponseWriter, r requestPayload) {
	data, err := msgpack.Marshal(&pingResponsePayload{
		CMD:  "response",
		To:   r.ReqID,
		Body: "pong",
	})
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func router(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload
	if err := msgpack.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch payload.CMD {
	case "ping":
		pingHandler(w, payload)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) Listen() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", router)

	srv := http.Server{
		Addr:    ":43111",
		Handler: mux,
	}
	s.srv = &srv

	println("file server listening...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
