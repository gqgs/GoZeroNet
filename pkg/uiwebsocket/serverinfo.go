package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/info"
)

type (
	serverInfoResponse struct {
		required
		Result info.Server `json:"result"`
	}
)

func (w *uiWebsocket) serverInfo(message Message) error {
	return w.conn.WriteJSON(serverInfoResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		info.ServerInfo(false),
	})
}
