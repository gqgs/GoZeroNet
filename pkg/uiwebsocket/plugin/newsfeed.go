package plugin

import (
	"encoding/json"
	"fmt"
)

type newsFeedPlugin struct{}

func NewNewsFeedPlugin() Plugin {
	return &newsFeedPlugin{}
}

func (newsFeedPlugin) Name() string {
	return "Newsfeed"
}

func (newsFeedPlugin) Handles(cmd string) bool {
	switch cmd {
	case "feedFollow", "feedListFollow", "feedQuery":
		return true
	default:
		return false
	}
}

func (p newsFeedPlugin) Handle(w pluginWriter, cmd string, to, id int64, message []byte) error {
	switch cmd {
	case "feedQuery":
		return p.feedQuery(w, id, message)
	default:
		reply := fmt.Sprintf(`{"error": "cmd not implemented, "to": %d, "id", %d}`, to, id)
		w.Write([]byte(reply))
		return fmt.Errorf("TODO: implement me: %s", cmd)
	}
}

type (
	feedQueryRequest struct {
		CMD          string          `json:"cmd"`
		ID           int64           `json:"id"`
		Params       feedQueryParams `json:"params"`
		WrapperNonce string          `json:"wrapper_nonce"`
	}
	feedQueryParams struct {
		DayLimit int `json:"day_limit"`
		Limit    int `json:"int"`
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

func (newsFeedPlugin) feedQuery(w pluginWriter, id int64, message []byte) error {
	request := new(feedQueryRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(feedQueryResponse{
		CMD: "response",
		ID:  id,
		To:  request.ID,
		Result: feedQueryResult{
			Rows:  make([]string, 0),
			Stats: make([]string, 0),
		},
	})
}
