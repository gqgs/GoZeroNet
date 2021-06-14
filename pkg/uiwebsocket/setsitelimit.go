package uiwebsocket

import (
	"encoding/json"
	"errors"
)

type (
	siteSetLimitRequest struct {
		CMD    string `json:"cmd"`
		ID     int64  `json:"id"`
		Params []int  `json:"params"`
	}

	siteSetLimitResponse struct {
		CMD    string             `json:"cmd"`
		ID     int64              `json:"id"`
		To     int64              `json:"to"`
		Result siteSetLimitResult `json:"result"`
	}

	siteSetLimitResult string
)

func (w *uiWebsocket) siteSetLimit(rawMessage []byte, message Message) error {
	request := new(siteSetLimitRequest)
	if err := json.Unmarshal(rawMessage, request); err != nil {
		return err
	}

	if len(request.Params) == 0 {
		return errors.New("missing required parameter")
	}

	// TODO: admin only
	if err := w.site.SetSiteLimit(request.Params[0], w.currentUser); err != nil {
		return err
	}

	return w.conn.WriteJSON(siteSetLimitResponse{
		CMD:    "response",
		To:     message.ID,
		ID:     w.ID(),
		Result: "ok",
	})
}
