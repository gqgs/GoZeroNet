package fileserver

import (
	"net"

	"github.com/gqgs/go-zeronet/pkg/lib/log"
)

type conn struct {
	net.Conn
}

// Creates and returns a new connection to the address.
// If the returned error is nil the client must close the
// connection after using it.
func NewConnection(addr string) (net.Conn, error) {
	netConn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &conn{
		Conn: netConn,
	}, nil
}

func NewDebugConnection(addr string) (net.Conn, error) {
	conn, err := NewConnection(addr)
	if err != nil {
		return nil, err
	}
	return debugConn{
		Conn: conn,
		log:  log.New("connection"),
	}, nil
}

type debugConn struct {
	net.Conn
	log log.Logger
}

func (c debugConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	c.log.WithField("op", "Write").Debug(n, string(b))
	return
}

func (c debugConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.log.WithField("op", "Read").Debug(n, string(b))
	return
}
