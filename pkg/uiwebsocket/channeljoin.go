package uiwebsocket

import (
	"encoding/json"
)

type (
	channelJoinRequest struct {
		required
		Params json.RawMessage `json:"params"`
	}
	channelJoinParams struct {
		Channels []string `json:"channels"`
	}

	channelJoinResponse struct {
		required
		Result channelJoinResult `json:"result"`
	}

	channelJoinResult string
)

func (w *uiWebsocket) channelJoin(rawMessage []byte, message Message) error {
	payload := new(channelJoinRequest)
	if err := jsonUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	var params channelJoinParams
	if err := jsonUnmarshal(payload.Params, &params); err != nil {
		if err := jsonUnmarshal(payload.Params, &params.Channels); err != nil {
			return err
		}
	}

	w.channelsMutex.Lock()
	for _, channel := range params.Channels {
		w.channels[channel] = struct{}{}
	}
	w.channelsMutex.Unlock()

	return w.conn.WriteJSON(channelJoinResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
