package file

import (
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_GetFile(t *testing.T) {
	srv, err := NewServer(config.RandomIPv4Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Shutdown()
	go srv.Listen()

	conn, err := NewConnection(srv.addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	const site = "site"
	const innerPath = "innerPath"
	resp, err := GetFile(conn, site, innerPath)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "response", resp.CMD)
	assert.Equal(t, 1, resp.To)
}
