package uiwebsocket

type (
	serverErrorsRequest struct {
		CMD          string             `json:"cmd"`
		ID           int64              `json:"id"`
		Params       serverErrorsParams `json:"params"`
		WrapperNonce string             `json:"wrapper_nonce"`
	}
	serverErrorsParams map[string]string

	serverErrorsResponse struct {
		CMD    string             `json:"cmd"`
		ID     int64              `json:"id"`
		To     int64              `json:"to"`
		Result serverErrorsResult `json:"result"`
	}

	serverErrorsResult []string
)

func (w *uiWebsocket) serverErrors(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(serverErrorsResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: make(serverErrorsResult, 0),
	})
}
