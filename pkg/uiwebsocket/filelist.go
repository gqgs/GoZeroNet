package uiwebsocket

import "encoding/json"

type (
	fileListRequest struct {
		CMD    string         `json:"cmd"`
		ID     int64          `json:"id"`
		Params fileListParams `json:"params"`
	}
	fileListParams struct {
		InnerPath string `json:"inner_path"`
	}

	fileListResponse struct {
		CMD    string         `json:"cmd"`
		ID     int64          `json:"id"`
		To     int64          `json:"to"`
		Result fileListResult `json:"result"`
	}

	fileListResult []string
)

func (w *uiWebsocket) fileList(rawMessage []byte, message Message) error {
	payload := new(fileListRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	list, err := w.site.ListFiles(payload.Params.InnerPath)
	if err != nil {
		return err
	}

	return w.conn.WriteJSON(fileListResponse{
		CMD:    "response",
		ID:     w.ID(),
		To:     message.ID,
		Result: list,
	})
}
