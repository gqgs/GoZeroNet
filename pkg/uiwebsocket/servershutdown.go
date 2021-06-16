package uiwebsocket

import (
	"os"
)

type (
	serverShutdownRequest struct {
		CMD          string               `json:"cmd"`
		ID           int64                `json:"id"`
		Params       serverShutdownParams `json:"params"`
		WrapperNonce string               `json:"wrapper_nonce"`
	}
	serverShutdownParams map[string]string

	serverShutdownResponse struct {
		CMD    string               `json:"cmd"`
		ID     int64                `json:"id"`
		To     int64                `json:"to"`
		Result serverShutdownResult `json:"result"`
	}

	serverShutdownResult map[string]interface{}
)

func (w *uiWebsocket) serverShutdown(rawMessage []byte, message Message) error {
	// TODO: admin only
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return process.Signal(os.Interrupt)
}
