package uiwebsocket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

type (
	fileGetRequest struct {
		required
		Params json.RawMessage `json:"params"`
	}
	fileGetParams struct {
		InnerPath string `json:"inner_path"`
		Required  bool   `json:"required"`
		Format    string `json:"format"`
		Timeout   uint   `json:"timeout"`
	}

	fileGetResponse struct {
		required
		Result *string `json:"result"`
	}
)

func (w *uiWebsocket) fileGet(rawMessage []byte, message Message) error {
	payload := new(fileGetRequest)
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	params := new(fileGetParams)
	if err := serialize.JSONUnmarshal(payload.Params, params); err != nil {
		if err := serialize.JSONUnmarshal(payload.Params, &params.InnerPath); err != nil {
			return err
		}
	}

	var writer io.Writer
	reader := new(strings.Builder)
	if params.Format == "base64" {
		writer = base64.NewEncoder(base64.StdEncoding, reader)
		defer writer.(io.Closer).Close()
	} else {
		writer = reader
	}

	var timeout time.Duration
	if params.Required {
		timeout = config.FileNeedDeadline
		if params.Timeout > 0 {
			timeout = time.Duration(params.Timeout * uint(time.Second))
		}
	}
	ctx, cancel := context.WithTimeout(w.ctx, timeout)
	defer cancel()

	if err := w.site.ReadFile(ctx, strings.TrimSuffix(params.InnerPath, "|all"), writer); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	if closer, ok := writer.(io.Closer); ok {
		closer.Close()
	}

	var result *string
	if r := reader.String(); len(r) > 0 {
		result = &r
	}

	return w.conn.WriteJSON(fileGetResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		result,
	})
}
