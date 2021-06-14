package uiwebsocket

type (
	userGetGlobalSettingsResponse struct {
		CMD    string                      `json:"cmd"`
		ID     int64                       `json:"id"`
		To     int64                       `json:"to"`
		Result userGetGlobalSettingsResult `json:"result"`
	}

	userGetGlobalSettingsResult map[string]interface{}
)

func (w *uiWebsocket) userGetGlobalSettings(rawMessage []byte, message Message) error {
	settings := w.currentUser.GlobalSettings()
	if len(settings) == 0 {
		settings = make(userGetGlobalSettingsResult)
	}
	return w.conn.WriteJSON(userGetGlobalSettingsResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     w.ID(),
		Result: settings,
	})
}
