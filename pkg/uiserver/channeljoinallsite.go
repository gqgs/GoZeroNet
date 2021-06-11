package uiserver

type (
	channelJoinAllsiteRequest struct {
		CMD          string                   `json:"cmd"`
		ID           int                      `json:"id"`
		Params       channelJoinAllsiteParams `json:"params"`
		WrapperNonce string                   `json:"wrapper_nonce"`
	}
	channelJoinAllsiteParams struct {
		Channel string `json:"channel"`
	}

	channelJoinAllsiteResponse struct {
		CMD    string                   `json:"cmd"`
		ID     int                      `json:"id"`
		To     int                      `json:"to"`
		Result channelJoinAllsiteResult `json:"result"`
	}

	channelJoinAllsiteResult string
)

func (w *uiWebsocket) channelJoinAllsite(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(channelJoinAllsiteResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     message.ID,
		Result: "ok",
	})
}
