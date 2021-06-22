package plugin

import (
	"encoding/json"

	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/site"
)

type optionalManager struct {
	ID IDFunc
}

func NewOptionalManager(idFunc IDFunc) Plugin {
	return &optionalManager{
		ID: idFunc,
	}
}

func (*optionalManager) Name() string {
	return "OptionalManager"
}

func (*optionalManager) Description() string {
	return "Manage optional content"
}

func (o *optionalManager) Handler(cmd string) (HandlerFunc, bool) {
	switch cmd {
	case "optionalLimitStats":
		return o.optionalLimitStats, true
	case "optionalFileInfo":
		return o.optionalFileInfo, true
	// case "optionalFileList", "optionalFilePin", "optionalFileUnpin",
	// 	"optionalFileDelete", "optionalLimitSet", "optionalHelpList",
	// 	"optionalHelp", "optionalHelpRemove", "optionalHelpAll":
	default:
		return nil, false
	}
}

type (
	optionalLimitStatsRequest struct {
		CMD    string                   `json:"cmd"`
		ID     int64                    `json:"id"`
		Params optionalLimitStatsParams `json:"params"`
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
		Used  int    `json:"used"`
	}
)

func (o *optionalManager) optionalLimitStats(w pluginWriter, site *site.Site, message []byte) error {
	request := new(optionalLimitStatsRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(optionalLimitStatsResponse{
		CMD: "response",
		ID:  o.ID(),
		To:  request.ID,
		Result: optionalLimitStatsResult{
			Free:  540246016,
			Limit: "10%",
		},
	})
}

type (
	optionalFileInfoRequest struct {
		CMD    string                 `json:"cmd"`
		ID     int64                  `json:"id"`
		Params optionalFileInfoParams `json:"params"`
	}
	optionalFileInfoParams struct {
		InnerPath string `json:"inner_path"`
	}

	optionalFileInfoResponse struct {
		CMD    string          `json:"cmd"`
		ID     int64           `json:"id"`
		To     int64           `json:"to"`
		Result *event.FileInfo `json:"result"`
	}
)

func (o *optionalManager) optionalFileInfo(w pluginWriter, site *site.Site, message []byte) error {
	request := new(optionalFileInfoRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}

	info, err := site.FileInfo(request.Params.InnerPath)
	if err != nil && err != database.ErrFileNotFound {
		return err
	}

	return w.WriteJSON(optionalFileInfoResponse{
		CMD:    "response",
		ID:     o.ID(),
		To:     request.ID,
		Result: info,
	})
}
