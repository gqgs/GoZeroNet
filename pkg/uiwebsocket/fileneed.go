package uiwebsocket

import "encoding/json"

type (
	fileNeedRequest struct {
		required
		Params fileNeedParams `json:"params"`
	}
	fileNeedParams struct {
		InnerPath string `json:"inner_path"`
	}

	fileNeedResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) fileNeed(rawMessage []byte, message Message) error {
	payload := new(fileNeedRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	if err := w.site.FileNeed(payload.Params.InnerPath); err != nil {
		return err
	}

	return w.conn.WriteJSON(fileNeedResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
