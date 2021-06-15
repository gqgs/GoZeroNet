package uiwebsocket

import "github.com/gqgs/go-zeronet/pkg/user"

type (
	userGetGlobalSettingsResponse struct {
		CMD    string              `json:"cmd"`
		ID     int64               `json:"id"`
		To     int64               `json:"to"`
		Result user.GlobalSettings `json:"result"`
	}

	userGetGlobalSettingsResult map[string]user.GlobalSettings
)

func (w *uiWebsocket) userGetGlobalSettings(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(userGetGlobalSettingsResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     w.ID(),
		Result: w.site.User().GlobalSettings(),
	})
}
