package uiserver

import "github.com/gqgs/go-zeronet/pkg/announcer"

type (
	announcerStatsRequest struct {
		CMD          string               `json:"cmd"`
		ID           int                  `json:"id"`
		Params       announcerStatsParams `json:"params"`
		WrapperNonce string               `json:"wrapper_nonce"`
	}
	announcerStatsParams map[string]string

	announcerStatsResponse struct {
		CMD    string               `json:"cmd"`
		ID     int                  `json:"id"`
		To     int                  `json:"to"`
		Result announcerStatsResult `json:"result"`
	}

	announcerStatsResult map[string]announcer.Stats
)

func (w *uiWebsocket) announcerStats(message []byte, id int) {
	err := w.conn.WriteJSON(announcerStatsResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     id,
		Result: announcer.GetStats(),
	})
	w.log.IfError(err)
}
