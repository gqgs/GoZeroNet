package file

import (
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_Handshake(t *testing.T) {
	srv, err := NewServer(config.RandomIPv4Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Shutdown()
	go srv.Listen()

	clientFileServer, err := NewServer(config.RandomIPv4Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer clientFileServer.Shutdown()

	conn, err := NewConnection(srv.addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	resp, err := Handshake(conn, srv.addr, clientFileServer)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "response", resp.CMD)
	assert.Equal(t, 1, resp.To)
	assert.Equal(t, "tls-rsa", resp.Crypt)
	assert.Equal(t, srv.port, resp.FileserverPort)
	assert.Equal(t, config.Protocol, resp.Protocol)
	assert.Equal(t, config.PortOpened, resp.PortOpened)
	assert.Equal(t, srv.peerID, resp.PeerID)
	assert.Equal(t, config.Rev, resp.Rev)
	assert.Equal(t, config.UseBinType, resp.UseBinType)
	assert.Equal(t, config.Version, resp.Version)
}
