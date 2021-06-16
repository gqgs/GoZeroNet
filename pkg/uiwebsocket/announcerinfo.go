package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	announcerInfoRequest struct {
		CMD    string              `json:"cmd"`
		ID     int64               `json:"id"`
		Params announcerInfoParams `json:"params"`
	}
	announcerInfoParams map[string]string

	announcerInfoResponse struct {
		CMD    string              `json:"cmd"`
		ID     int64               `json:"id"`
		To     int64               `json:"to"`
		Result announcerInfoResult `json:"result"`
	}

	announcerInfoResult struct {
		Address string                          `json:"address"`
		Stats   map[string]*site.AnnouncerStats `json:"stats"`
	}
)

func (w *uiWebsocket) announcerInfo(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(announcerInfoResponse{
		CMD: "response",
		ID:  w.ID(),
		To:  message.ID,
		Result: announcerInfoResult{
			Address: w.site.Address(),
			Stats:   w.site.AnnouncerStats(),
		},
	})
}
