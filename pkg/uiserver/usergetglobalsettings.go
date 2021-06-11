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

func (w *uiWebsocket) userGetGlobalSettings(message []byte, id int) {
	err := w.conn.WriteJSON(userGetGlobalSettingsResponse{
		CMD:    "response",
		To:     id,
		ID:     w.reqID,
		Result: make(userGetGlobalSettingsResult),
	})
	w.log.IfError(err)
}
