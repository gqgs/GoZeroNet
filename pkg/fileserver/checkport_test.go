package fileserver

import (
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/stretchr/testify/assert"
)

func Test_CheckPort(t *testing.T) {
	srv, err := NewServer(config.RandomIPv4Addr)
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

	t.Run("open", func(t *testing.T) {
		resp, err := CheckPort(conn, srv.port)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "response", resp.CMD)
		assert.Equal(t, 1, resp.To)
		assert.Equal(t, "open", resp.Status)
		assert.Equal(t, conn.LocalAddr().String(), resp.IPExternal)
	})

	t.Run("closed", func(t *testing.T) {
		resp, err := CheckPort(conn, 1337)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "response", resp.CMD)
		assert.Equal(t, 1, resp.To)
		assert.Equal(t, "closed", resp.Status)
		assert.Equal(t, conn.LocalAddr().String(), resp.IPExternal)
	})
}
