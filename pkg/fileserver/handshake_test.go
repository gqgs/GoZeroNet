package fileserver

import (
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/stretchr/testify/assert"
)

func Test_Handshake(t *testing.T) {
	pubsubManagerSrv := pubsub.NewManager()
	contentDBSrv, err := database.NewContentDatabase()
	if err != nil {
		t.Fatal(err)
	}
	defer contentDBSrv.Close()

	pubsubManagerClient := pubsub.NewManager()
	contentDBClient, err := database.NewContentDatabase()
	if err != nil {
		t.Fatal(err)
	}
	defer contentDBClient.Close()

	srv, err := NewServer(config.RandomIPv4Addr, contentDBSrv, pubsubManagerSrv)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Shutdown()
	go srv.Listen()

	clientFileServer, err := NewServer(config.RandomIPv4Addr, contentDBClient, pubsubManagerClient)
	if err != nil {
		t.Fatal(err)
	}
	defer clientFileServer.Shutdown()

	conn, err := connection.NewConnection(srv.addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	resp, err := Handshake(conn, srv.addr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "response", resp.CMD)
	assert.Equal(t, "", resp.Crypt)
	assert.Equal(t, config.Protocol, resp.Protocol)
	assert.Equal(t, config.PortOpened, resp.PortOpened)
	assert.Equal(t, config.Rev, resp.Rev)
	assert.Equal(t, config.UseBinType, resp.UseBinType)
	assert.Equal(t, config.Version, resp.Version)
}
