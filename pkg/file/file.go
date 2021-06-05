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

// reads only the necessary to decode the cmd
func decodeCmd(decoder *msgpack.Decoder) (string, error) {
	var buffer bytes.Buffer
	decoder.Reset(io.TeeReader(decoder.Buffered(), &buffer))
	query, err := decoder.Query("cmd")
	if err != nil {
		return "", err
	}
	if len(query) != 1 {
		return "", errors.New("file: invalid cmd value")
	}

	cmd, ok := query[0].(string)
	if !ok {
		return "", errors.New("file: invalid cmd value")
	}
	decoder.Reset(io.MultiReader(&buffer, decoder.Buffered()))
	return cmd, nil
}

func decode(reader io.Reader) (interface{}, error) {
	decoder := msgpack.NewDecoder(reader)
	cmd, err := decodeCmd(decoder)
	if err != nil {
		return nil, err
	}

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
