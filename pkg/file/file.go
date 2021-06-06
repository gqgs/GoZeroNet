package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

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
	id := random.PeerID()
	return &server{
		peerID: id,
		log:    log.New("fileserver").WithField("peerid", id),
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	if s == nil || s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *server) Listen() {
	s.log.Infof("listening at %s", config.FileServer.Addr())

	l, err := net.Listen("tcp", config.FileServer.Addr())
	if err != nil {
		s.log.Fatal(err)
	}

	for {
		// TODO: should check if error implements net.Error interface
		// and try again if the error is temporary
		conn, err := l.Accept()
		if err != nil {
			s.log.Error(err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *server) handleConn(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(config.ReadDeadline))
	if err := s.route(conn, conn); err != nil {
		s.log.Error(err)
	}
}

func (s *server) route(w io.Writer, r io.Reader) error {
	s.log.Debug("new request")
	i, err := decode(r)
	if err != nil {
		return err
	}

	switch req := i.(type) {
	case pingRequest:
		return s.pingHandler(w, req)
	case handshakeRequest:
		return s.handshakeHandler(w, req)
	case getFileRequest:
		return s.getFileHandler(w, req)
	default:
		return errors.New("file: invalid command")
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
