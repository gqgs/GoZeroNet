package plugin

import (
	"github.com/gqgs/go-zeronet/pkg/site"
)

type newsFeedPlugin struct {
	ID IDFunc
}

func NewNewsFeed(idFunc IDFunc) Plugin {
	return &newsFeedPlugin{
		ID: idFunc,
	}
}

func (*newsFeedPlugin) Name() string {
	return "Newsfeed"
}

func (*newsFeedPlugin) Description() string {
	return "Feeds from SQL queries"
}

func (n *newsFeedPlugin) Handler(cmd string) (HandlerFunc, bool) {
	switch cmd {
	case "feedQuery":
		return n.feedQuery, true
	// case "feedListFollow", "feedFollow"
	default:
		return nil, false
	}
}

type (
	feedQueryRequest struct {
		required
		Params feedQueryParams `json:"params"`
	}
	feedQueryParams struct {
		DayLimit int `json:"day_limit"`
		Limit    int `json:"int"`
	}

	feedQueryResponse struct {
		required
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

func (n *newsFeedPlugin) feedQuery(w pluginWriter, site *site.Site, message []byte) error {
	request := new(feedQueryRequest)
	if err := jsonUnmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(feedQueryResponse{
		required{
			CMD: "response",
			ID:  n.ID(),
			To:  request.ID,
		},
		feedQueryResult{
			Rows:  make([]string, 0),
			Stats: make([]string, 0),
		},
	})
}
