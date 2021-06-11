package uiserver

import "sync/atomic"

type (
	userGetGlobalSettingsResponse struct {
		CMD    string                      `json:"cmd"`
		ID     int64                       `json:"id"`
		To     int64                       `json:"to"`
		Result userGetGlobalSettingsResult `json:"result"`
	}

	userGetGlobalSettingsResult map[string]string
)

func (w *uiWebsocket) userGetGlobalSettings(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(userGetGlobalSettingsResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     atomic.AddInt64(&w.reqID, 1),
		Result: make(userGetGlobalSettingsResult),
	})
}
