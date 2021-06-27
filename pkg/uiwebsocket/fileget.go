package uiwebsocket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
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

	var timeout time.Duration
	if payload.Params.Required {
		timeout = config.FileNeedDeadline
		if payload.Params.Timeout > 0 {
			timeout = time.Duration(payload.Params.Timeout * uint(time.Second))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := w.site.ReadFile(ctx, strings.TrimSuffix(payload.Params.InnerPath, "|all"), writer); err != nil {
		if !errors.Is(err, database.ErrFileNotFound) {
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
