package fileserver

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/spf13/cast"
	"github.com/vmihailenco/msgpack/v5"
)

var counter = safe.Counter()

// Implements the protocol specified at:
// https://zeronet.io/docs/help_zeronet/network_protocol/
//
// Every message is encoded with MessagePack.
// Every request has 3 parameters: `cmd`, `req_id` and `params`.
type server struct {
	l             net.Listener
	log           log.Logger
	addr          string // host:port
	host          string
	port          int
	contentDB     database.ContentDatabase
	pubsubManager pubsub.Manager
}

type Server interface{}

func NewServer(addr string, contentDB database.ContentDatabase, pubsubManager pubsub.Manager) (*server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	// Can be different from `addr` if the port was chosen by the server.
	chosenAddr := l.Addr().String()
	host, portString, _ := net.SplitHostPort(chosenAddr)
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}

	config.FileServerHost = host
	config.FileServerPort = port

	return &server{
		addr:          chosenAddr,
		port:          port,
		host:          host,
		l:             l,
		log:           log.New("fileserver"),
		contentDB:     contentDB,
		pubsubManager: pubsubManager,
	}, nil
}

func (s *server) Shutdown() error {
	if s == nil || s.l == nil {
		return nil
	}
	return s.l.Close()
}

func (s *server) Listen() {
	s.log.Infof("listening at http://%s", s.addr)
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
		if config.Debug {
			go s.handleConn(connection.NewTraceConn(conn))
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *server) handleConn(conn net.Conn) {
	logger := s.log.
		WithField("local", conn.LocalAddr()).
		WithField("remote", conn.RemoteAddr())
	logger.Debug("new connection")

	defer conn.Close()

	for {
		if err := s.route(conn); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error(err)
			return
		}
	}
}

func (s *server) route(conn net.Conn) error {
	decoder := msgpack.NewDecoder(conn)
	cmd, err := decodeKey(decoder, "cmd")
	if err != nil {
		return err
	}

	switch cast.ToString(cmd) {
	case "ping":
		return s.pingHandler(conn, decoder)
	case "handshake":
		return s.handshakeHandler(conn, decoder)
	case "getFile":
		return s.getFileHandler(conn, decoder)
	case "streamFile":
		return s.streamFileHandler(conn, decoder)
	case "checkport":
		return s.checkPortHandler(conn, decoder)
	case "pex":
		return s.pexHandler(conn, decoder)
	case "listModified":
		return s.listModifiedHandler(conn, decoder)
	case "update":
		return s.updateHandler(conn, decoder)
	case "findHashIds":
		return s.findHashIDsHandler(conn, decoder)
	case "getHashfield":
		return s.getHashfieldHandler(conn, decoder)
	case "setHashfield":
		return s.setHashfieldHandler(conn, decoder)
	case "getPiecefields":
		return s.getPieceFieldsHandler(conn, decoder)
	case "setPiecefields":
		return s.setPieceFieldsHandler(conn, decoder)
	default:
		s.log.WithField("cmd", cmd).Warn("unknown request")
		return s.unknownHandler(conn, decoder)
	}
}

type requestDecoder interface {
	Reset(io.Reader)
	Decode(v interface{}) error
	Buffered() io.Reader
	Query(query string) ([]interface{}, error)
}

func decodeKey(decoder requestDecoder, key string) (interface{}, error) {
	var buffer bytes.Buffer
	decoder.Reset(io.TeeReader(decoder.Buffered(), &buffer))
	defer decoder.Reset(io.MultiReader(&buffer, decoder.Buffered()))
	query, err := decoder.Query(key)
	if err != nil {
		return "", err
	}
	if len(query) != 1 {
		return "", errors.New("file: invalid cmd value")
	}
	return query[0], nil
}

func peerID(conn net.Conn) string {
	if peer, ok := conn.(interface {
		ID() string
	}); ok {
		return peer.ID()
	}
	return random.PeerID()
}
