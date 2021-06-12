package uiwebsocket

import "sync/atomic"

type (
	channelInfoRequest struct {
		CMD    string            `json:"cmd"`
		ID     int64             `json:"id"`
		Params channelInfoParams `json:"params"`
	}
	channelInfoParams struct {
		Channels []string `json:"channels"`
	}

	channelInfoResponse struct {
		CMD    string            `json:"cmd"`
		ID     int64             `json:"id"`
		To     int64             `json:"to"`
		Result channelInfoResult `json:"result"`
	}

	channelInfoResult string
)

func (w *uiWebsocket) channelInfo(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(channelInfoResponse{
		CMD:    "response",
		ID:     atomic.AddInt64(&w.reqID, 1),
		To:     message.ID,
		Result: "ok",
	})
}
