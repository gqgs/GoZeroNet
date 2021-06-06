package file

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func Test_GetFile(t *testing.T) {
	srv := NewServer()
	go srv.Listen()
	defer srv.Shutdown(context.Background())

	body := &getFileRequest{
		CMD:   "getFile",
		ReqID: 1,
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
	var decoded getFileResponse
	assert.NoError(t, msgpack.NewDecoder(resp.Body).Decode(&decoded))
	assert.Equal(t, "response", decoded.CMD)
	assert.Equal(t, body.ReqID, decoded.To)
	assert.Equal(t, "TODO", decoded.Body)
	assert.Equal(t, "TODO", decoded.Location)
}
