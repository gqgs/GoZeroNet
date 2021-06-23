package uiwebsocket

type (
	userGetSettingsResponse struct {
		required
		Result userGetSettingsResult `json:"result"`
	}

	userGetSettingsResult map[string]interface{}
)

func (w *uiWebsocket) userGetSettings(message Message) error {
	settings := w.site.User().SiteSettings(w.site.Address())
	if len(settings) == 0 {
		settings = make(userGetSettingsResult)
	}
	return w.conn.WriteJSON(userGetSettingsResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		settings,
	})
}
