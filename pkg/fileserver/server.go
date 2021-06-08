package fileserver

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/spf13/cast"
	"github.com/vmihailenco/msgpack/v5"
)

// Implements the protocol specified at:
// https://zeronet.io/docs/help_zeronet/network_protocol/
//
// Every message is encoded with MessagePack.
// Every request has 3 parameters: `cmd`, `req_id` and `params`.
type server struct {
	l      net.Listener
	log    log.Logger
	peerID string
	addr   string // host:port
	host   string
	port   int
}

func NewServer(addr string) (*server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	// Can be different from `addr` if the port was chosen by the server.
	chosenAddr := l.Addr().String()
	hostString, portString, _ := net.SplitHostPort(chosenAddr)
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}

	return &server{
		peerID: random.PeerID(),
		addr:   chosenAddr,
		port:   port,
		host:   hostString,
		l:      l,
		log:    log.New("fileserver"),
	}, nil
}

func (s *server) Shutdown() error {
	if s == nil || s.l == nil {
		return nil
	}
	return s.l.Close()
}

func (s *server) Listen() {
	s.log.Infof("listening at %s", s.addr)
	for {
		conn, err := s.l.Accept()
		switch e := err.(type) {
		case *net.OpError:
			if e.Temporary() {
				s.log.Warn(err)
				continue
			}
			switch e.Err.Error() {
			case "use of closed network connection":
				return
			}
			s.log.Error(err)
			return
		case nil:
		default:
			s.log.Error(e)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *server) handleConn(conn net.Conn) {
	s.log.
		WithField("local", conn.LocalAddr()).
		WithField("remote", conn.RemoteAddr().String()).
		Debug("new connection")

	defer conn.Close()
	conn.SetDeadline(time.Now().Add(config.Deadline))

	for {
		if err := s.route(conn); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			s.log.Error(err)
			return
		}
	}
}

func (s *server) route(conn net.Conn) error {
	s.log.Debug("new request")
	i, err := decode(conn)
	if err != nil {
		return err
	}

	switch r := i.(type) {
	case pingRequest:
		return pingHandler(conn, r)
	case handshakeRequest:
		return handshakeHandler(conn, r, s)
	case getFileRequest:
		return getFileHandler(conn, r)
	case streamFileRequest:
		return streamFileHandler(conn, r)
	case checkPortRequest:
		return checkPortHandler(conn, r, s)
	case unknownRequest:
		return unknownHandler(conn, r)
	default:
		return errors.New("file: invalid command")
	}
}

func decode(reader io.Reader) (interface{}, error) {
	decoder := msgpack.NewDecoder(reader)
	cmd, err := decodeKey(decoder, "cmd")
	if err != nil {
		return nil, err
	}

	switch cast.ToString(cmd) {
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
	case "streamFile":
		var payload streamFileRequest
		err := decoder.Decode(&payload)
		return payload, err
	case "checkport":
		var payload checkPortRequest
		err := decoder.Decode(&payload)
		return payload, err
	default:
		reqID, err := decodeKey(decoder, "req_id")
		if err != nil {
			return nil, err
		}
		return unknownRequest{
			ReqID: cast.ToInt(reqID),
		}, nil
	}
}

// Reads only the necessary to decode the key.
func decodeKey(decoder *msgpack.Decoder, key string) (interface{}, error) {
	var buffer bytes.Buffer
	decoder.Reset(io.TeeReader(decoder.Buffered(), &buffer))
	query, err := decoder.Query(key)
	if err != nil {
		return "", err
	}
	if len(query) != 1 {
		return "", errors.New("file: invalid cmd value")
	}
	decoder.Reset(io.MultiReader(&buffer, decoder.Buffered()))
	return query[0], nil
}
