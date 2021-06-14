package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	siteListRequest struct {
		CMD          string         `json:"cmd"`
		ID           int64          `json:"id"`
		Params       siteListParams `json:"params"`
		WrapperNonce string         `json:"wrapper_nonce"`
	}
	siteListParams struct {
		ConnectionSites bool `json:"connecting_sites"`
	}

	siteListResponse struct {
		CMD    string       `json:"cmd"`
		ID     int64        `json:"id"`
		To     int64        `json:"to"`
		Result []*site.Info `json:"result"`
	}

	siteListResult string
)

func (w *uiWebsocket) siteList(rawMessage []byte, message Message) error {
	info, err := w.site.Info(w.currentUser)
	if err != nil {
		return err
	}

	return w.conn.WriteJSON(siteListResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: []*site.Info{info},
	})
}
