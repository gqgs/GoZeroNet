package uiwebsocket

import (
	"encoding/json"
	"errors"
)

type (
	dbQueryRequest struct {
		CMD    string   `json:"cmd"`
		ID     int64    `json:"id"`
		Params []string `json:"params"`
	}

	dbQueryResponse struct {
		CMD    string                   `json:"cmd"`
		ID     int64                    `json:"id"`
		To     int64                    `json:"to"`
		Result []map[string]interface{} `json:"result"`
	}

	dbQueryResult []string
)

func (w *uiWebsocket) dbQuery(rawMessage []byte, message Message) error {
	payload := new(dbQueryRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
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
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: result,
	})
}
