package file

import (
	"bytes"
	"context"
	"errors"
	"io"
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
	if s == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func decode(reader io.Reader) (interface{}, error) {
	// read only the necessary to decode the cmd
	var buffer bytes.Buffer
	decoder := msgpack.NewDecoder(io.TeeReader(reader, &buffer))
	query, err := decoder.Query("cmd")
	if err != nil {
		return nil, err
	}
	if len(query) != 1 {
		return nil, errors.New("file: bad cmd value")
	}

	cmd, ok := query[0].(string)
	if !ok {
		return nil, errors.New("file: bad cmd value")
	}

	// now reset and read the payload
	decoder.Reset(io.MultiReader(&buffer, decoder.Buffered()))

	switch cmd {
	case "ping":
		var payload pingRequest
		err := decoder.Decode(&payload)
		return payload, err
	case "handshake":
		var payload handshakeRequest
		err := decoder.Decode(&payload)
		return payload, err
	}

	return nil, errors.New("file: invalid payload type")
}

func router(w http.ResponseWriter, r *http.Request) {
	i, err := decode(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	switch v := i.(type) {
	case pingRequest:
		pingHandler(w, v)
	case handshakeRequest:
		handshakeHandler(w, v)
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
