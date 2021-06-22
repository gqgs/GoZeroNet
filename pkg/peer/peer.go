package peer

import (
	"fmt"
	"net"

	"github.com/gqgs/go-zeronet/pkg/connection"
)

type peer struct {
	net.Conn
	connected bool
	addr      string
}

type Peer interface {
	fmt.Stringer
	net.Conn
	Connect() error
}

func (p *peer) Connect() error {
	if p.connected {
		return nil
	}

	conn, err := connection.NewConnection(p.addr)
	if err != nil {
		return err
	}

	p.connected = true
	p.Conn = conn
	return nil
}

func (p *peer) Close() error {
	if p.Conn != nil {
		return p.Conn.Close()
	}
	p.connected = false
	return nil
}

func (p *peer) String() string {
	return p.addr
}

func NewPeer(addr string) *peer {
	return &peer{
		addr: addr,
	}
}
