package uiserver

type (
	serverErrorsRequest struct {
		CMD          string             `json:"cmd"`
		ID           int                `json:"id"`
		Params       serverErrorsParams `json:"params"`
		WrapperNonce string             `json:"wrapper_nonce"`
	}
	serverErrorsParams map[string]string

	serverErrorsResponse struct {
		CMD    string             `json:"cmd"`
		ID     int                `json:"id"`
		To     int                `json:"to"`
		Result serverErrorsResult `json:"result"`
	}

	serverErrorsResult []string
)

func (w *uiWebsocket) serverErrors(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(serverErrorsResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     message.ID,
		Result: make(serverErrorsResult, 0),
	})
}
