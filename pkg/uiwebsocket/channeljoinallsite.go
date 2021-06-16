package uiwebsocket

import (
	"encoding/json"
)

type (
	channelJoinAllsiteRequest struct {
		CMD    string                   `json:"cmd"`
		ID     int64                    `json:"id"`
		Params channelJoinAllsiteParams `json:"params"`
	}
	channelJoinAllsiteParams struct {
		Channel string `json:"channel"`
	}

	channelJoinAllsiteResponse struct {
		CMD    string                   `json:"cmd"`
		ID     int64                    `json:"id"`
		To     int64                    `json:"to"`
		Result channelJoinAllsiteResult `json:"result"`
	}

	channelJoinAllsiteResult string
)

func (w *uiWebsocket) channelJoinAllsite(rawMessage []byte, message Message) error {
	payload := new(channelJoinAllsiteRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	w.channelsMutex.Lock()
	w.channels[payload.Params.Channel] = struct{}{}
	w.channelsMutex.Unlock()

	w.allChannels = true

	return w.conn.WriteJSON(channelJoinAllsiteResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: "ok",
	})
}
