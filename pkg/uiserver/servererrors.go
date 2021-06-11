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

func (w *uiWebsocket) serverErrors(message []byte, id int) {
	err := w.conn.WriteJSON(serverErrorsResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     id,
		Result: make(serverErrorsResult, 0),
	})
	w.log.IfError(err)
}
