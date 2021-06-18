package connection

import (
	"net"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
)

type conn struct {
	net.Conn
}

// Creates and returns a new connection to the address.
// If the returned error is nil the client must close the
// connection after using it.
func NewConnection(addr string) (net.Conn, error) {
	netConn, err := net.DialTimeout("tcp", addr, config.ConnectionDeadline)
	if err != nil {
		return nil, err
	}

	if config.Debug {
		return traceConn{
			Conn: netConn,
			log:  log.New("connection"),
		}, nil
	}

	return &conn{
		Conn: netConn,
	}, nil
}

type traceConn struct {
	net.Conn
	log log.Logger
}

func (c traceConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	c.log.WithField("op", "Write").Trace(n, string(b))
	return
}

func (c traceConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.log.WithField("op", "Read").Trace(n, string(b))
	return
}
