package uiwebsocket

import "sync/atomic"

type (
	siteLimitRequest struct {
		CMD    string          `json:"cmd"`
		ID     int64           `json:"id"`
		Params siteLimitParams `json:"params"`
	}
	siteLimitParams struct {
		Channels []string `json:"channels"`
	}

	siteLimitResponse struct {
		CMD    string          `json:"cmd"`
		ID     int64           `json:"id"`
		To     int64           `json:"to"`
		Result siteLimitResult `json:"result"`
	}

	siteLimitResult string
)

func (w *uiWebsocket) siteLimit(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(siteLimitResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     atomic.AddInt64(&w.reqID, 1),
		Result: "ok",
	})
}
