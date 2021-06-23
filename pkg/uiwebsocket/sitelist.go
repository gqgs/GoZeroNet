package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	siteListResponse struct {
		required
		Result []*site.Info `json:"result"`
	}
)

func (w *uiWebsocket) siteList(rawMessage []byte, message Message) error {
	siteList, err := w.siteManager.SiteList()
	if err != nil {
		return err
	}

	return w.conn.WriteJSON(siteListResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		siteList,
	})
}
