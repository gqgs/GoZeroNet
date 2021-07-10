package uiwebsocket

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

type (
	fileWriteRequest struct {
		required
		Params []string `json:"params"`
	}

	fileWriteResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) fileWrite(rawMessage []byte, message Message) error {
	payload := new(fileWriteRequest)
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	if len(payload.Params) < 2 {
		return errors.New("invalid request")
	}

	innerPath := payload.Params[0]
	contentBase64 := payload.Params[1]

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(contentBase64))
	if err := w.site.FileWrite(innerPath, reader); err != nil {
		return err
	}

	return w.conn.WriteJSON(fileWriteResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
