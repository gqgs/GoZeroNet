package uiserver

type (
	optionalLimitStatsRequest struct {
		CMD          string                   `json:"cmd"`
		ID           int                      `json:"id"`
		Params       optionalLimitStatsParams `json:"params"`
		WrapperNonce string                   `json:"wrapper_nonce"`
	}
	optionalLimitStatsParams map[string]string

	optionalLimitStatsResponse struct {
		CMD    string                   `json:"cmd"`
		ID     int                      `json:"id"`
		To     int                      `json:"to"`
		Result optionalLimitStatsResult `json:"result"`
	}

	optionalLimitStatsResult struct {
		Free  int    `json:"free"`
		Limit string `json:"limit"`
		Used  int    `json:"usd"`
	}
)

func (w *uiWebsocket) optionalLimitStats(message []byte, id int) {
	err := w.conn.WriteJSON(optionalLimitStatsResponse{
		CMD: "response",
		ID:  w.reqID,
		To:  id,
		Result: optionalLimitStatsResult{
			Free:  540246016,
			Limit: "10%",
		},
	})
	w.log.IfError(err)
}
