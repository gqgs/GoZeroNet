package uiserver

type (
	optionalLimitStatsRequest struct {
		CMD          string                   `json:"cmd"`
		ID           int64                    `json:"id"`
		Params       optionalLimitStatsParams `json:"params"`
		WrapperNonce string                   `json:"wrapper_nonce"`
	}
	optionalLimitStatsParams map[string]string

	optionalLimitStatsResponse struct {
		CMD    string                   `json:"cmd"`
		ID     int64                    `json:"id"`
		To     int64                    `json:"to"`
		Result optionalLimitStatsResult `json:"result"`
	}

	optionalLimitStatsResult struct {
		Free  int    `json:"free"`
		Limit string `json:"limit"`
		Used  int    `json:"usd"`
	}
)

func (w *uiWebsocket) optionalLimitStats(rawMessage []byte, message Message) error {
	return w.conn.WriteJSON(optionalLimitStatsResponse{
		CMD: "response",
		ID:  w.reqID,
		To:  message.ID,
		Result: optionalLimitStatsResult{
			Free:  540246016,
			Limit: "10%",
		},
	})
}
