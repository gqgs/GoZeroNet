package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	siteInfoResponse struct {
		required
		Result *site.Info `json:"result"`
	}
)

func (w *uiWebsocket) siteInfo(message Message) error {
	info, err := w.site.Info()
	if err != nil {
		return err
	}

	return w.conn.WriteJSON(siteInfoResponse{
		required{
			CMD: "response",
			To:  message.ID,
			ID:  w.ID(),
		},
		info,
	})
}
