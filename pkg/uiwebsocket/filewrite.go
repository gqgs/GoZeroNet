package uiwebsocket

import (
	"encoding/base64"
	"strings"
)

type (
	fileWriteRequest struct {
		required
		Params fileWriteParams `json:"params"`
	}
	fileWriteParams struct {
		InnerPath     string `json:"inner_path"`
		ContentBase64 string `json:"content_base64"`
	}

	fileWriteResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) fileWrite(rawMessage []byte, message Message) error {
	payload := new(fileWriteRequest)
	if err := jsonUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(payload.Params.ContentBase64))
	if err := w.site.FileWrite(payload.Params.InnerPath, reader); err != nil {
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
