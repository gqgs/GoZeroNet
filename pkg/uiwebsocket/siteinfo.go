package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	siteInfoRequest struct {
		CMD    string         `json:"cmd"`
		ID     int64          `json:"id"`
		Params siteInfoParams `json:"params"`
	}
	siteInfoParams map[struct{}]struct{}

	siteInfoResponse struct {
		CMD    string     `json:"cmd"`
		ID     int64      `json:"id"`
		To     int64      `json:"to"`
		Result *site.Info `json:"result"`
	}
)

func (w *uiWebsocket) siteInfo(rawMessage []byte, message Message) error {
	info, err := w.site.Info()
	if err != nil {
		return err
	}

	return w.conn.WriteJSON(siteInfoResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     w.ID(),
		Result: info,
	})
}
