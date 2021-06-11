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

func (w *uiWebsocket) channelInfo(message []byte, id int) {
	err := w.conn.WriteJSON(channelInfoResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     id,
		Result: "ok",
	})
	w.log.IfError(err)
}
