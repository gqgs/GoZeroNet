package uiserver

import "sync/atomic"

type (
	feedQueryRequest struct {
		CMD          string          `json:"cmd"`
		ID           int64           `json:"id"`
		Params       feedQueryParams `json:"params"`
		WrapperNonce string          `json:"wrapper_nonce"`
	}
	feedQueryParams struct {
		DayLimit int `json:"day_limit`
		Limit    int `json:"int`
	}

	feedQueryResponse struct {
		CMD    string          `json:"cmd"`
		ID     int64           `json:"id"`
		To     int64           `json:"to"`
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

func (w *uiWebsocket) feedQuery(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(feedQueryResponse{
		CMD: "response",
		ID:  atomic.AddInt64(&w.reqID, 1),
		To:  message.ID,
		Result: feedQueryResult{
			Rows:  make([]string, 0),
			Stats: make([]string, 0),
		},
	})
}
