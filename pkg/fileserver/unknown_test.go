package fileserver

import (
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func Test_Unknown(t *testing.T) {
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

	encoded, err := msgpack.Marshal(&pingRequest{
		CMD:    "pong",
		ReqID:  1,
		Params: make(map[string]struct{}),
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err = conn.Write(encoded); err != nil {
		t.Fatal(err)
	}

	resp := new(unknownResponse)
	if err = msgpack.NewDecoder(conn).Decode(resp); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "response", resp.CMD)
	assert.Equal(t, "Unknown cmd", resp.Error)
}
