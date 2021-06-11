package uiserver

type (
	feedQueryRequest struct {
		CMD          string          `json:"cmd"`
		ID           int             `json:"id"`
		Params       feedQueryParams `json:"params"`
		WrapperNonce string          `json:"wrapper_nonce"`
	}
	feedQueryParams struct {
		DayLimit int `json:"day_limit`
		Limit    int `json:"int`
	}

	feedQueryResponse struct {
		CMD    string          `json:"cmd"`
		ID     int             `json:"id"`
		To     int             `json:"to"`
		Result feedQueryResult `json:"result"`
	}

	feedQueryResult struct {
		Num   int      `json:"num"`
		Rows  []string `json:"rows"`
		Sites int      `json:"sites"`
		Stats []string `json:"stats"`
		Taken int      `json:"taken"`
	}
)

func (w *uiWebsocket) feedQuery(message []byte, id int) {
	err := w.conn.WriteJSON(feedQueryResponse{
		CMD: "response",
		ID:  w.reqID,
		To:  id,
		Result: feedQueryResult{
			Rows:  make([]string, 0),
			Stats: make([]string, 0),
		},
	})
	w.log.IfError(err)
}
