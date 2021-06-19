package uiwebsocket

import "encoding/json"

type (
	userSetSettingsRequest struct {
		CMD    string                `json:"cmd"`
		ID     int64                 `json:"id"`
		Params userSetSettingsParams `json:"params"`
	}
	userSetSettingsParams struct {
		Settings map[string]interface{} `json:"settings"`
	}

	userSetSettingsResponse struct {
		CMD    string                `json:"cmd"`
		ID     int64                 `json:"id"`
		To     int64                 `json:"to"`
		Result userSetSettingsResult `json:"result"`
	}

	userSetSettingsResult string
)

func (w *uiWebsocket) userSetSettings(rawMessage []byte, message Message) error {
	payload := new(userSetSettingsRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	if err := w.site.User().SetSiteSettings(w.site.Address(), payload.Params.Settings); err != nil {
		return err
	}

	return w.conn.WriteJSON(userSetSettingsResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: "ok",
	})
}
