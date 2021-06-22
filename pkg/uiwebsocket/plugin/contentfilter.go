package plugin

import (
	"encoding/json"

	"github.com/gqgs/go-zeronet/pkg/site"
)

type contentFilter struct {
	ID IDFunc
}

func NewContentFilter(idFunc IDFunc) Plugin {
	return &contentFilter{
		ID: idFunc,
	}
}

func (*contentFilter) Name() string {
	return "ContentFilter"
}

func (*contentFilter) Description() string {
	return "Manage site and user block list"
}

func (c *contentFilter) Handler(cmd string) (HandlerFunc, bool) {
	switch cmd {
	case "filterIncludeList":
		return c.filterIncludeList, true
	// case "filterIncludeAdd", "filterIncludeRemove":
	default:
		return nil, false
	}
}

type (
	filterIncludeListRequest struct {
		CMD    string                  `json:"cmd"`
		ID     int64                   `json:"id"`
		Params filterIncludeListParams `json:"params"`
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

func (c *contentFilter) filterIncludeList(w pluginWriter, site *site.Site, message []byte) error {
	request := new(filterIncludeListRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(filterIncludeListResponse{
		CMD:    "response",
		ID:     c.ID(),
		To:     request.ID,
		Result: make(filterIncludeListResult, 0),
	})
}
