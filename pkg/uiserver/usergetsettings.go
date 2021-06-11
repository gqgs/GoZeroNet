package uiserver

type (
	userGetSettingsRequest struct {
		CMD          string                `json:"cmd"`
		ID           int                   `json:"id"`
		Params       userGetSettingsParams `json:"params"`
		WrapperNonce string                `json:"wrapper_nonce"`
	}
	userGetSettingsParams map[string]string

	userGetSettingsResponse struct {
		CMD    string                `json:"cmd"`
		ID     int                   `json:"id"`
		To     int                   `json:"to"`
		Result userGetSettingsResult `json:"result"`
	}

	userGetSettingsResult map[string]string
)

func (w *uiWebsocket) userGetSettings(message []byte, id int) {
	err := w.conn.WriteJSON(userGetSettingsResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     id,
		Result: make(userGetSettingsResult),
	})
	w.log.IfError(err)
}
