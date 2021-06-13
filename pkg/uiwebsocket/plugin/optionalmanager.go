package plugin

import (
	"encoding/json"
	"fmt"
)

type optionalManager struct{}

func (optionalManager) Name() string {
	return "OptionalManager"
}

func NewOptionalManager() Plugin {
	return &optionalManager{}
}

func (optionalManager) Description() string {
	return "Manage optional content"
}

func (optionalManager) Handles(cmd string) bool {
	switch cmd {
	case "optionalFileList", "optionalFileInfo", "optionalFilePin", "optionalFileUnpin",
		"optionalFileDelete", "optionalLimitStats", "optionalLimitSet", "optionalHelpList",
		"optionalHelp", "optionalHelpRemove", "optionalHelpAll":
		return true
	default:
		return false
	}
}

func (o optionalManager) Handle(w pluginWriter, cmd string, to, id int64, message []byte) error {
	switch cmd {
	case "optionalLimitStats":
		return o.optionalLimitStats(w, id, message)
	default:
		reply := fmt.Sprintf(`{"error": "not implemented: %q", "to": %d, "id", %d}`, cmd, to, id)
		w.Write([]byte(reply))
		return fmt.Errorf("TODO: implement me: %s", cmd)
	}
}

type (
	optionalLimitStatsRequest struct {
		CMD          string                   `json:"cmd"`
		ID           int64                    `json:"id"`
		Params       optionalLimitStatsParams `json:"params"`
		WrapperNonce string                   `json:"wrapper_nonce"`
	}
	optionalLimitStatsParams []string

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

func (optionalManager) optionalLimitStats(w pluginWriter, id int64, message []byte) error {
	request := new(optionalLimitStatsRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(optionalLimitStatsResponse{
		CMD: "response",
		ID:  id,
		To:  request.ID,
		Result: optionalLimitStatsResult{
			Free:  540246016,
			Limit: "10%",
		},
	})
}
