package uiwebsocket

import (
	"encoding/json"
)

type (
	channelJoinRequest struct {
		CMD    string          `json:"cmd"`
		ID     int64           `json:"id"`
		Params json.RawMessage `json:"params"`
	}
	channelJoinParams struct {
		Channels []string `json:"channels"`
	}

	channelJoinResponse struct {
		CMD    string            `json:"cmd"`
		ID     int64             `json:"id"`
		To     int64             `json:"to"`
		Result channelJoinResult `json:"result"`
	}

	channelJoinResult string
)

func (w *uiWebsocket) channelJoin(rawMessage []byte, message Message) error {
	payload := new(channelJoinRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	var params channelJoinParams
	if err := json.Unmarshal(payload.Params, &params); err != nil {
		if err := json.Unmarshal(payload.Params, &params.Channels); err != nil {
			return err
		}
	}

	w.channelsMutex.Lock()
	for _, channel := range params.Channels {
		w.channels[channel] = struct{}{}
	}
	w.channelsMutex.Unlock()

	return w.conn.WriteJSON(channelJoinResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: "ok",
	})
}
