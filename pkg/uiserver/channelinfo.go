package uiserver

type (
	channelInfoRequest struct {
		CMD    string            `json:"cmd"`
		ID     int               `json:"id"`
		Params channelInfoParams `json:"params"`
	}
	channelInfoParams struct {
		Channels []string `json:"channels"`
	}

	channelInfoResponse struct {
		CMD    string            `json:"cmd"`
		ID     int               `json:"id"`
		To     int               `json:"to"`
		Result channelInfoResult `json:"result"`
	}

	channelInfoResult string
)

func (w *uiWebsocket) channelInfo(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(channelInfoResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     message.ID,
		Result: "ok",
	})
}
