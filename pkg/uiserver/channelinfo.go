package uiserver

type (
	ChannelInfoRequest struct {
		CMD    string            `json:"cmd"`
		ID     int               `json:"id"`
		Params ChannelInfoParams `json:"params"`
	}
	ChannelInfoParams struct {
		Channels []string `json:"channels"`
	}

	ChannelInfoResponse struct {
		CMD    string            `json:"cmd"`
		ID     int               `json:"id"`
		To     int               `json:"to"`
		Result ChannelInfoResult `json:"result"`
	}

	ChannelInfoResult string
)

func (w *uiWebsocket) channelInfo(message []byte, id int) {
	err := w.conn.WriteJSON(ChannelInfoResponse{
		CMD:    "response",
		To:     id,
		Result: "ok",
	})
	w.log.IfError(err)
}
