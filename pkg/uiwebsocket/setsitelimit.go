package uiwebsocket

import (
	"errors"

	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

type (
	siteSetLimitRequest struct {
		required
		Params []int `json:"params"`
	}

	siteSetLimitResponse struct {
		required
		Result siteSetLimitResult `json:"result"`
	}

	siteSetLimitResult string
)

func (w *uiWebsocket) siteSetLimit(rawMessage []byte, message Message) error {
	request := new(siteSetLimitRequest)
	if err := serialize.JSONUnmarshal(rawMessage, request); err != nil {
		return err
	}

	if len(request.Params) == 0 {
		return errors.New("missing required parameter")
	}

	if err := w.site.SetSiteLimit(request.Params[0]); err != nil {
		return err
	}

	return w.conn.WriteJSON(siteSetLimitResponse{
		required{
			CMD: "response",
			To:  message.ID,
			ID:  w.ID(),
		},
		"ok",
	})
}
