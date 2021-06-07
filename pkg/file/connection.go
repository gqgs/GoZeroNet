package file

import (
	"net"

	"github.com/gqgs/go-zeronet/pkg/lib/log"
)

type conn struct {
	net.Conn
	log  log.Logger
	addr string
}

func NewConnection(addr string) (*conn, error) {
	netConn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &conn{
		Conn: netConn,
		log:  log.New("connection"),
		addr: addr,
	}, nil
}
