package uiwebsocket

type (
	userGetSettingsRequest struct {
		CMD          string                `json:"cmd"`
		ID           int64                 `json:"id"`
		Params       userGetSettingsParams `json:"params"`
		WrapperNonce string                `json:"wrapper_nonce"`
	}
	userGetSettingsParams map[string]string

	userGetSettingsResponse struct {
		CMD    string                `json:"cmd"`
		ID     int64                 `json:"id"`
		To     int64                 `json:"to"`
		Result userGetSettingsResult `json:"result"`
	}

	userGetSettingsResult map[string]string
)

func (w *uiWebsocket) userGetSettings(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(userGetSettingsResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: make(userGetSettingsResult),
	})
}
