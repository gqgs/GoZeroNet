package fileserver

import (
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/stretchr/testify/assert"
)

func Test_Pex(t *testing.T) {
	srv, err := NewServer(config.RandomIPv4Addr, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Shutdown()
	go srv.Listen()

	conn, err := connection.NewConnection(srv.addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	const site = "site"
	resp, err := Pex(conn, site, 5)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "response", resp.CMD)
	assert.Len(t, resp.Peers, 0)
	assert.Len(t, resp.PeersOnion, 0)
}
