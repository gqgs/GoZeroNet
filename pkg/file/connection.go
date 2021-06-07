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

// Creates and returns a new connection to the address.
// If the returned error is nil the client must close the
// connection after using it.
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
