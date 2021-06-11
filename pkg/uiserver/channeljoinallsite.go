package uiserver

import "sync/atomic"

type (
	channelJoinAllsiteRequest struct {
		CMD          string                   `json:"cmd"`
		ID           int64                    `json:"id"`
		Params       channelJoinAllsiteParams `json:"params"`
		WrapperNonce string                   `json:"wrapper_nonce"`
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
	return w.conn.WriteJSON(channelJoinAllsiteResponse{
		CMD:    "response",
		ID:     atomic.AddInt64(&w.reqID, 1),
		To:     message.ID,
		Result: "ok",
	})
}
