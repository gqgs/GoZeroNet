package uiserver

//go:generate go run github.com/mailru/easyjson/easyjson -all

type (
	UserGetGlobalSettingsResponse struct {
		CMD    string                      `json:"cmd"`
		ID     int                         `json:"id"`
		To     int                         `json:"to"`
		Result UserGetGlobalSettingsResult `json:"result"`
	}

	UserGetGlobalSettingsResult map[string]string
)

func (w *uiWebsocket) UserGetGlobalSettings(message []byte, id int) {
	err := w.conn.WriteJSON(UserGetGlobalSettingsResponse{
		CMD:    "response",
		To:     id,
		Result: make(UserGetGlobalSettingsResult),
	})
	w.log.IfError(err)
}
