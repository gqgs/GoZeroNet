package plugin

import (
	"encoding/json"
	"fmt"
)

type contentFilter struct{}

func NewContentFilter() Plugin {
	return &contentFilter{}
}

func (contentFilter) Name() string {
	return "ContentFilter"
}

func (contentFilter) Description() string {
	return "Manage site and user block list"
}

func (contentFilter) Handles(cmd string) bool {
	switch cmd {
	case "filterIncludeAdd", "filterIncludeRemove", "filterIncludeList":
		return true
	default:
		return false
	}
}

func (c contentFilter) Handle(w pluginWriter, cmd string, to, id int64, message []byte) error {
	switch cmd {
	case "filterIncludeList":
		return c.filterIncludeList(w, id, message)
	default:
		w.WriteJSON(errorMsg{
			Msg: "not implemented",
			To:  to,
			ID:  id,
		})
		return fmt.Errorf("TODO: implement me: %s", cmd)
	}
}

type (
	filterIncludeListRequest struct {
		CMD          string                  `json:"cmd"`
		ID           int64                   `json:"id"`
		Params       filterIncludeListParams `json:"params"`
		WrapperNonce string                  `json:"wrapper_nonce"`
	}
	filterIncludeListParams struct {
		AllSites bool `json:"all_sites"`
		Filters  bool `json:"filters"`
	}

	filterIncludeListResponse struct {
		CMD    string                  `json:"cmd"`
		ID     int64                   `json:"id"`
		To     int64                   `json:"to"`
		Result filterIncludeListResult `json:"result"`
	}

	filterIncludeListResult []string
)

func (contentFilter) filterIncludeList(w pluginWriter, id int64, message []byte) error {
	request := new(filterIncludeListRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(filterIncludeListResponse{
		CMD:    "response",
		ID:     id,
		To:     request.ID,
		Result: make(filterIncludeListResult, 0),
	})
}
