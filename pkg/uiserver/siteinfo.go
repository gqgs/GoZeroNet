package uiserver

import "github.com/gqgs/go-zeronet/pkg/site"

type (
	siteInfoRequest struct {
		CMD    string         `json:"cmd"`
		ID     int            `json:"id"`
		Params siteInfoParams `json:"params"`
	}
	siteInfoParams map[struct{}]struct{}

	siteInfoResponse struct {
		CMD    string    `json:"cmd"`
		ID     int       `json:"id"`
		To     int       `json:"to"`
		Result site.Info `json:"result"`
	}
)

func (w *uiWebsocket) siteInfo(message []byte, id int) {
	err := w.conn.WriteJSON(siteInfoResponse{
		CMD:    "response",
		To:     id,
		ID:     w.reqID,
		Result: site.GetInfo(w.siteManager),
	})
	w.log.IfError(err)
}
