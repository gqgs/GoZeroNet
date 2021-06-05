package file

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func Test_Handshake(t *testing.T) {
	srv := Server{}
	go srv.Listen()
	defer srv.Shutdown(context.Background())

	body := &handshakeRequest{
		CMD:   "handshake",
		ReqID: 1,
		HandshakeParams: handshakeParams{
			Rev:            2092,
			PortOpened:     false,
			FileserverPort: 43111,
			Protocol:       "v2",
		},
	}
	encoded, err := msgpack.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:43111/", bytes.NewReader(encoded))
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
	assert.Equal(t, 2092, decoded.Rev)
	assert.Equal(t, false, decoded.PortOpened)
	assert.Equal(t, 43111, decoded.FileserverPort)
	assert.Equal(t, "v2", decoded.Protocol)
}
