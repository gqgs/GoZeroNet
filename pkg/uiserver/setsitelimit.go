package uiserver

type (
	siteLimitRequest struct {
		CMD    string          `json:"cmd"`
		ID     int             `json:"id"`
		Params siteLimitParams `json:"params"`
	}
	siteLimitParams struct {
		Channels []string `json:"channels"`
	}

	siteLimitResponse struct {
		CMD    string          `json:"cmd"`
		ID     int             `json:"id"`
		To     int             `json:"to"`
		Result siteLimitResult `json:"result"`
	}

	siteLimitResult string
)

func (w *uiWebsocket) siteLimit(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(siteLimitResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     w.reqID,
		Result: "ok",
	})
}
