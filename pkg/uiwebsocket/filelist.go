package uiwebsocket

type (
	fileListRequest struct {
		required
		Params fileListParams `json:"params"`
	}
	fileListParams struct {
		InnerPath string `json:"inner_path"`
	}

	fileListResponse struct {
		required
		Result fileListResult `json:"result"`
	}

	fileListResult []string
)

func (w *uiWebsocket) fileList(rawMessage []byte, message Message) error {
	payload := new(fileListRequest)
	if err := jsonUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	list, err := w.site.ListFiles(payload.Params.InnerPath)
	if err != nil {
		return err
	}

	return w.conn.WriteJSON(fileListResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		list,
	})
}
