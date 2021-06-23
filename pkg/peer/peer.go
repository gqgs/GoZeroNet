package peer

import (
	"fmt"
	"net"

	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
)

type peer struct {
	net.Conn
	connected bool
	addr      string
	id        string
}

type Peer interface {
	fmt.Stringer
	net.Conn
	Connect() error
	ID() string
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

func (p *peer) ID() string {
	return p.id
}

func NewPeer(addr string) *peer {
	return &peer{
		id:   random.PeerID(),
		addr: addr,
	}
}
