package uiwebsocket

import (
	"errors"
)

func (w *uiWebsocket) response(rawMessage []byte, message Message) error {
	w.waitingMutex.Lock()
	defer w.waitingMutex.Unlock()

	fn, ok := w.waitingResponses[message.To]
	if !ok {
		return errors.New("waiting response not found")
	}
	delete(w.waitingResponses, message.To)

	return fn(rawMessage)
}
