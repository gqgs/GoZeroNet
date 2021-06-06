package file

import (
	"context"
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
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

	resp, err := client.Ping(srv.Addr())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "response", resp.CMD)
	assert.Equal(t, 1, resp.To)
	assert.Equal(t, "Pong!", resp.Body)
}
