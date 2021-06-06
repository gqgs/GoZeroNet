package file

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func Test_Handshake(t *testing.T) {
	srv := NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go srv.Listen(ctx)

	body := &handshakeRequest{
		CMD:   "handshake",
		ReqID: 1,
		Params: handshakeParams{
			Crypt:          "tls-rsa",
			CryptSupported: []string{"tls-rsa"},
			FileserverPort: config.FileServer.Port,
			Protocol:       config.Protocol,
			PortOpened:     config.PortOpened,
			PeerID:         random.PeerID(),
			Rev:            config.Rev,
			UseBinType:     config.UseBinType,
			Version:        config.Version,
		},
	}
	encoded, err := msgpack.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, testURL(), bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var decoded handshakeResponse
	assert.NoError(t, msgpack.NewDecoder(resp.Body).Decode(&decoded))
	assert.Equal(t, "response", decoded.CMD)
	assert.Equal(t, body.ReqID, decoded.To)
	assert.Equal(t, body.Params.Rev, decoded.Rev)
	assert.Equal(t, body.Params.Version, decoded.Version)
	assert.Equal(t, body.Params.PortOpened, decoded.PortOpened)
	assert.Equal(t, body.Params.FileserverPort, decoded.FileserverPort)
	assert.Equal(t, body.Params.Protocol, decoded.Protocol)
	assert.Equal(t, srv.peerID, decoded.PeerID)
}
