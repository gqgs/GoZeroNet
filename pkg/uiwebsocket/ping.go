package uiwebsocket

type (
	pingResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) ping(message Message) error {
	return w.conn.WriteJSON(pingResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"Pong!",
	})
}
