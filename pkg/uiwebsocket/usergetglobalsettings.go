package uiwebsocket

import "github.com/gqgs/go-zeronet/pkg/user"

type (
	userGetGlobalSettingsResponse struct {
		required
		Result user.GlobalSettings `json:"result"`
	}
)

func (w *uiWebsocket) userGetGlobalSettings(message Message) error {
	return w.conn.WriteJSON(userGetGlobalSettingsResponse{
		required{
			CMD: "response",
			To:  message.ID,
			ID:  w.ID(),
		},
		w.site.User().GlobalSettings(),
	})
}
