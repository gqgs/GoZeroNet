package uiwebsocket

type (
	fileDeleteRequest struct {
		required
		Params fileDeleteParams `json:"params"`
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

	if err := w.site.FileDelete(payload.Params.InnerPath); err != nil {
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
