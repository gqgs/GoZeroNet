package uiserver

import "github.com/gqgs/go-zeronet/pkg/info"

type (
	serverInfoRequest struct {
		CMD          string           `json:"cmd"`
		ID           int              `json:"id"`
		Params       serverInfoParams `json:"params"`
		WrapperNonce string           `json:"wrapper_nonce"`
	}
	serverInfoParams map[string]string

	serverInfoResponse struct {
		CMD    string      `json:"cmd"`
		ID     int         `json:"id"`
		To     int         `json:"to"`
		Result info.Server `json:"result"`
	}
)

func (w *uiWebsocket) serverInfo(message []byte, id int) {
	err := w.conn.WriteJSON(serverInfoResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     id,
		Result: info.ServerInfo(),
	})
	w.log.IfError(err)
}
