package uiwebsocket

type (
	channelJoinAllsiteRequest struct {
		required
		Params channelJoinAllsiteParams `json:"params"`
	}
	channelJoinAllsiteParams struct {
		Channel string `json:"channel"`
	}

	channelJoinAllsiteResponse struct {
		required
		Result channelJoinAllsiteResult `json:"result"`
	}

	channelJoinAllsiteResult string
)

func (w *uiWebsocket) channelJoinAllsite(rawMessage []byte, message Message) error {
	payload := new(channelJoinAllsiteRequest)
	if err := jsonUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	w.channelsMutex.Lock()
	w.channels[payload.Params.Channel] = struct{}{}
	w.channelsMutex.Unlock()

	w.allChannels = true

	return w.conn.WriteJSON(channelJoinAllsiteResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
