package uiwebsocket

import (
	"os"
)

func (w *uiWebsocket) serverShutdown(rawMessage []byte, message Message) error {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return process.Signal(os.Interrupt)
}
