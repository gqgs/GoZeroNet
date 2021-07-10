package uiwebsocket

import (
	"errors"

	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

type (
	dbQueryRequest struct {
		required
		Params []string `json:"params"`
	}

	dbQueryResponse struct {
		required
		Result []map[string]interface{} `json:"result"`
	}
)

func (w *uiWebsocket) dbQuery(rawMessage []byte, message Message) error {
	payload := new(dbQueryRequest)
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	if len(payload.Params) == 0 {
		return errors.New("missing query")
	}

	result, err := w.site.Query(payload.Params[0])
	if err != nil {
		return err
	}

	return w.conn.WriteJSON(dbQueryResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		result,
	})
}
