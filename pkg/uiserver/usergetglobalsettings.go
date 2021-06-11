package uiserver

type (
	userGetGlobalSettingsResponse struct {
		CMD    string                      `json:"cmd"`
		ID     int                         `json:"id"`
		To     int                         `json:"to"`
		Result userGetGlobalSettingsResult `json:"result"`
	}

	userGetGlobalSettingsResult map[string]string
)

func (w *uiWebsocket) userGetGlobalSettings(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(userGetGlobalSettingsResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     w.reqID,
		Result: make(userGetGlobalSettingsResult),
	})
}
