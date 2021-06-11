package uiserver

import "github.com/gqgs/go-zeronet/pkg/site"

type (
	siteListRequest struct {
		CMD          string         `json:"cmd"`
		ID           int            `json:"id"`
		Params       siteListParams `json:"params"`
		WrapperNonce string         `json:"wrapper_nonce"`
	}
	siteListParams struct {
		ConnectionSites bool `json:"connecting_sites"`
	}

	siteListResponse struct {
		CMD    string      `json:"cmd"`
		ID     int         `json:"id"`
		To     int         `json:"to"`
		Result []site.Info `json:"result"`
	}

	siteListResult string
)

func (w *uiWebsocket) siteList(message []byte, id int) {
	err := w.conn.WriteJSON(siteListResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     id,
		Result: []site.Info{site.GetInfo(w.siteManager)},
	})
	w.log.IfError(err)
}
