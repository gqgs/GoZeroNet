package uiwebsocket

import "encoding/json"

type (
	fileDeleteRequest struct {
		required
		Params json.RawMessage `json:"params"`
	}
	fileDeleteParams struct {
		InnerPath string `json:"inner_path"`
	}

	fileDeleteResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) fileDelete(rawMessage []byte, message Message) error {
	payload := new(fileDeleteRequest)
	if err := jsonUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	params := new(fileDeleteParams)
	if err := jsonUnmarshal(payload.Params, params); err != nil {
		if err := jsonUnmarshal(payload.Params, &params.InnerPath); err != nil {
			return err
		}
	}

	if err := w.site.FileDelete(params.InnerPath); err != nil {
		return err
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
