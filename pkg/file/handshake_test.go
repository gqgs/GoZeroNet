package file

import (
	"context"
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_Handshake(t *testing.T) {
	srv, err := NewServer(config.RandomIPv4Addr)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go srv.Listen(ctx)

	client, err := NewServer(config.RandomIPv4Addr)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Handshake(srv.addr)
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
