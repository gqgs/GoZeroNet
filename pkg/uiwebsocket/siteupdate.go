package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

type (
	siteUpdateRequest struct {
		required
		Params siteUpdateParams `json:"params"`
	}

	siteUpdateParams struct {
		Address string `json:"address"`
	}
	siteUpdateResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) siteUpdate(rawMessage []byte, message Message) error {
	payload := new(siteUpdateRequest)
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	event.BroadcastSiteUpdate(payload.Params.Address, w.pubsubManager, &event.SiteUpdate{})

	return w.conn.WriteJSON(siteUpdateResponse{
		required{
			CMD: "response",
			To:  message.ID,
			ID:  w.ID(),
		},
		"Updated",
	})
}
