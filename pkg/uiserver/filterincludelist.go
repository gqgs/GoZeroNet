package uiserver

type (
	filterIncludeListRequest struct {
		CMD          string                  `json:"cmd"`
		ID           int                     `json:"id"`
		Params       filterIncludeListParams `json:"params"`
		WrapperNonce string                  `json:"wrapper_nonce"`
	}
	filterIncludeListParams struct {
		AllSites bool `json:"all_sites"`
		Filters  bool `json:"filters"`
	}

	filterIncludeListResponse struct {
		CMD    string                  `json:"cmd"`
		ID     int                     `json:"id"`
		To     int                     `json:"to"`
		Result filterIncludeListResult `json:"result"`
	}

	filterIncludeListResult []string
)

func (w *uiWebsocket) filterIncludeList(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(filterIncludeListResponse{
		CMD:    "response",
		ID:     w.reqID,
		To:     message.ID,
		Result: make(filterIncludeListResult, 0),
	})
}
