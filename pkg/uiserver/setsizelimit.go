package uiserver

type (
	sizeLimitRequest struct {
		CMD    string          `json:"cmd"`
		ID     int             `json:"id"`
		Params sizeLimitParams `json:"params"`
	}
	sizeLimitParams struct {
		Channels []string `json:"channels"`
	}

	sizeLimitResponse struct {
		CMD    string          `json:"cmd"`
		ID     int             `json:"id"`
		To     int             `json:"to"`
		Result sizeLimitResult `json:"result"`
	}

	sizeLimitResult string
)

func (w *uiWebsocket) siteLimit(message []byte, id int) {
	err := w.conn.WriteJSON(sizeLimitResponse{
		CMD:    "response",
		To:     id,
		Result: "ok",
	})
	w.log.IfError(err)
}
