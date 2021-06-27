package fileserver

import (
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/stretchr/testify/assert"
)

func Test_Update(t *testing.T) {
	pubsubManager := pubsub.NewManager()
	contentDB, err := database.NewContentDatabase()
	if err != nil {
		t.Fatal(err)
	}
	defer contentDB.Close()

	srv, err := NewServer(config.RandomIPv4Addr, contentDB, pubsubManager)
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
	const innerPath = "innerPath"
	resp, err := Update(conn, site, innerPath)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "response", resp.CMD)
	assert.True(t, resp.Ok)
}
