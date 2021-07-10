package uiwebsocket

import (
	"strings"

	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

type (
	fileNeedRequest struct {
		required
		Params fileNeedParams `json:"params"`
	}
	fileNeedParams struct {
		InnerPath string `json:"inner_path"`
	}

	fileNeedResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) fileNeed(rawMessage []byte, message Message) error {
	payload := new(fileNeedRequest)
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	w.site.FileNeed(strings.TrimSuffix(payload.Params.InnerPath, "|all"))

	return w.conn.WriteJSON(fileNeedResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
