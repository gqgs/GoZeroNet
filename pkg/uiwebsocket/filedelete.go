package uiwebsocket

import "github.com/gqgs/go-zeronet/pkg/lib/serialize"

type (
	fileDeleteRequest struct {
		required
		Params []string `json:"params"`
	}

	fileDeleteResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) fileDelete(rawMessage []byte, message Message) error {
	payload := new(fileDeleteRequest)
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	for _, innerPath := range payload.Params {
		if err := w.site.FileDelete(innerPath); err != nil {
			return err
		}
	}

	return w.conn.WriteJSON(fileDeleteResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
