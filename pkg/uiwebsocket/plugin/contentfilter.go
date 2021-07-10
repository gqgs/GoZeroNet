package plugin

import (
	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
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
		required
		Params filterIncludeListParams `json:"params"`
	}
	filterIncludeListParams struct {
		AllSites bool `json:"all_sites"`
		Filters  bool `json:"filters"`
	}

	filterIncludeListResponse struct {
		required
		Result filterIncludeListResult `json:"result"`
	}

	filterIncludeListResult []string
)

func (c *contentFilter) filterIncludeList(w pluginWriter, site *site.Site, message []byte) error {
	request := new(filterIncludeListRequest)
	if err := serialize.JSONUnmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(filterIncludeListResponse{
		required{
			CMD: "response",
			ID:  c.ID(),
			To:  request.ID,
		},
		make(filterIncludeListResult, 0),
	})
}
