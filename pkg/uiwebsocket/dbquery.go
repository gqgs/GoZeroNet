package uiwebsocket

import (
	"encoding/json"
)

type (
	dbQueryRequest struct {
		CMD    string          `json:"cmd"`
		ID     int64           `json:"id"`
		Params json.RawMessage `json:"params"`
	}

	dbQueryResponse struct {
		CMD    string        `json:"cmd"`
		ID     int64         `json:"id"`
		To     int64         `json:"to"`
		Result dbQueryResult `json:"result"`
	}

	dbQueryResult []string
)

func (w *uiWebsocket) dbQuery(rawMessage []byte, message Message) error {
	payload := new(dbQueryRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	w.log.WithField("rawMessage", string(payload.Params)).Error("implement me")

	return w.conn.WriteJSON(dbQueryResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: make(dbQueryResult, 0),
	})
}
