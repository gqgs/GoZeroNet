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

func (w *uiWebsocket) announcerStats(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(announcerStatsResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     message.ID,
		Result: announcer.GetStats(),
	})
}
