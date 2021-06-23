package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	announcerStatsResponse struct {
		required
		Result announcerStatsResult `json:"result"`
	}

	announcerStatsResult map[string]*site.AnnouncerStats
)

func (w *uiWebsocket) announcerStats(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(announcerStatsResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		w.site.AnnouncerStats(),
	})
}
