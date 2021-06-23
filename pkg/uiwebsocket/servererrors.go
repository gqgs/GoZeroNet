package uiwebsocket

type (
	serverErrorsResponse struct {
		required
		Result serverErrorsResult `json:"result"`
	}

	serverErrorsResult []string
)

func (w *uiWebsocket) serverErrors(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(serverErrorsResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		make(serverErrorsResult, 0),
	})
}
