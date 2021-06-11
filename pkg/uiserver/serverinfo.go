package uiserver

import (
	"sync/atomic"

	"github.com/gqgs/go-zeronet/pkg/info"
)

type (
	serverInfoRequest struct {
		CMD          string           `json:"cmd"`
		ID           int64            `json:"id"`
		Params       serverInfoParams `json:"params"`
		WrapperNonce string           `json:"wrapper_nonce"`
	}
	serverInfoParams map[string]string

	serverInfoResponse struct {
		CMD    string      `json:"cmd"`
		ID     int64       `json:"id"`
		To     int64       `json:"to"`
		Result info.Server `json:"result"`
	}
)

func (w *uiWebsocket) serverInfo(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(serverInfoResponse{
		CMD:    "response",
		ID:     atomic.AddInt64(&w.reqID, 1),
		To:     message.ID,
		Result: info.ServerInfo(),
	})
}
