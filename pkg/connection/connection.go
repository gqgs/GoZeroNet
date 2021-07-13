package connection

import (
	"net"
	"strings"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"golang.org/x/net/proxy"
)

type conn struct {
	net.Conn
}

func (c *conn) Write(b []byte) (n int, err error) {
	c.Conn.SetWriteDeadline(time.Now().Add(config.ConnectionTimeout))
	return c.Conn.Write(b)
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.Conn.SetReadDeadline(time.Now().Add(config.ConnectionTimeout))
	return c.Conn.Read(b)
}

// NewConnection creates and returns a new connection to the address.
// If the returned error is nil the client must close the
// connection after using it.
func NewConnection(addr string) (net.Conn, error) {
	netConn, err := dial(addr, config.ConnectionTimeout)
	if err != nil {
		return nil, err
	}

	return Wrap(netConn), nil
}

func dial(addr string, timeout time.Duration) (net.Conn, error) {
	if strings.Contains(addr, ".onion") {
		dialer, err := proxy.SOCKS5("tcp", config.TorAddress, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return dialer.Dial("tcp", addr)
	}
	return net.DialTimeout("tcp", addr, timeout)
}

// Wrap wraps a net.Conn to return another net.Coon
// with properly configured deadlines and trace levels.
func Wrap(c net.Conn) net.Conn {
	if config.Debug {
		return newTraceConn(c)
	}
	return &conn{c}
}

func newTraceConn(conn net.Conn) *traceConn {
	return &traceConn{
		conn,
		log.New("connection"),
	}
}

type traceConn struct {
	net.Conn
	log log.Logger
}

func (c *traceConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	c.log.WithField("op", "Write").Trace(n, string(b))
	return
}

func (c *traceConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.log.WithField("op", "Read").Trace(n, string(b))
	return
}
