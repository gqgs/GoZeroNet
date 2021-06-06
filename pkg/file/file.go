package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/vmihailenco/msgpack/v5"
)

// Implements the protocol specified at:
// https://zeronet.io/docs/help_zeronet/network_protocol/
//
// Every message is encoded with MessagePack
// Every request has 3 parameters: `cmd`, `req_id` and `params`
type server struct {
	srv    *http.Server
	peerID string
	log    log.Logger
}

func NewServer() *server {
	mux := http.NewServeMux()
	id := random.PeerID()
	s := &server{
		srv: &http.Server{
			Addr:    config.FileServer.Addr(),
			Handler: mux,
		},
		peerID: id,
		log:    log.New("fileserver").WithField("peerid", id),
	}
	mux.Handle("/", s)
	return s
}

func (s *server) Shutdown(ctx context.Context) error {
	if s == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *server) Listen() {
	s.log.Infof("listening at %s", config.FileServer.Addr())
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Fatal(err)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i, err := decode(r.Body)
	if err != nil {
		s.log.Error(err)
		return
	}

	switch r := i.(type) {
	case pingRequest:
		s.pingHandler(w, r)
	case handshakeRequest:
		s.handshakeHandler(w, r)
	case getFileRequest:
		s.getFileHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
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
	case "getFile":
		var payload getFileRequest
		err := decoder.Decode(&payload)
		return payload, err
	default:
		return nil, fmt.Errorf("file: invalid payload type (%q)", cmd)
	}
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
