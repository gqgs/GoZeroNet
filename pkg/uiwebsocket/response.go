package uiwebsocket

import (
	"encoding/json"
	"errors"
)

type (
	responseRequest struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) response(rawMessage []byte) error {
	request := new(responseRequest)
	if err := json.Unmarshal(rawMessage, request); err != nil {
		return err
	}

	w.waitingMutex.Lock()
	defer w.waitingMutex.Unlock()

	fn, ok := w.waitingResponses[request.To]
	if !ok {
		return errors.New("waiting response not found")
	}
	delete(w.waitingResponses, request.To)

	return fn(request.Result)
}
