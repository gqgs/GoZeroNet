package uiwebsocket

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"strings"
)

type (
	fileGetRequest struct {
		required
		Params fileGetParams `json:"params"`
	}
	fileGetParams struct {
		InnerPath string `json:"inner_path"`
		Required  bool   `json:"required"`
		Format    string `json:"format"`
		Timeout   uint   `json:"timeout"`
	}

	fileGetResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) fileGet(rawMessage []byte, message Message) error {
	payload := new(fileGetRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	var writer io.Writer
	reader := new(strings.Builder)
	if payload.Params.Format == "base64" {
		writer = base64.NewEncoder(base64.StdEncoding, reader)
		defer writer.(io.Closer).Close()
	} else {
		writer = reader
	}

	if err := w.site.ReadFile(payload.Params.InnerPath, writer); err != nil {
		// TODO: download file with timeout if required
		if !os.IsNotExist(err) {
			return err
		}
		if payload.Params.Format != "base64" {
			writer.Write([]byte("{}"))
		}
	}

	return w.conn.WriteJSON(fileGetResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		reader.String(),
	})
}
