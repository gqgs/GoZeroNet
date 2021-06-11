package uiserver

import (
	"sync/atomic"

	"github.com/gqgs/go-zeronet/pkg/announcer"
)

type (
	announcerStatsRequest struct {
		CMD          string               `json:"cmd"`
		ID           int64                `json:"id"`
		Params       announcerStatsParams `json:"params"`
		WrapperNonce string               `json:"wrapper_nonce"`
	}
	announcerStatsParams map[string]string

	announcerStatsResponse struct {
		CMD    string               `json:"cmd"`
		ID     int64                `json:"id"`
		To     int64                `json:"to"`
		Result announcerStatsResult `json:"result"`
	}

	announcerStatsResult map[string]announcer.Stats
)

func (w *uiWebsocket) announcerStats(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(announcerStatsResponse{
		CMD:    "response",
		ID:     atomic.AddInt64(&w.reqID, 1),
		To:     message.ID,
		Result: announcer.GetStats(),
	})
}
