package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	announcerInfoResponse struct {
		required
		Result announcerInfoResult `json:"result"`
	}

	announcerInfoResult struct {
		Address string                          `json:"address"`
		Stats   map[string]*site.AnnouncerStats `json:"stats"`
	}
)

func (w *uiWebsocket) announcerInfo(message Message) error {
	return w.conn.WriteJSON(announcerInfoResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		announcerInfoResult{
			Address: w.site.Address(),
			Stats:   w.site.AnnouncerStats(),
		},
	})
}
