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

func (w *uiWebsocket) siteInfo(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(siteInfoResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     w.reqID,
		Result: site.GetInfo(w.siteManager),
	})
}