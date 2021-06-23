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
	case "optionalHelpList":
		return o.optionalHelpList, true
	// case "optionalFileList", "optionalFilePin", "optionalFileUnpin",
	// 	"optionalFileDelete", "optionalLimitSet", "optionalHelpList",
	// 	"optionalHelp", "optionalHelpRemove", "optionalHelpAll":
	default:
		return nil, false
	}
}

type (
	optionalLimitStatsRequest struct {
		required
		Params optionalLimitStatsParams `json:"params"`
	}
	optionalLimitStatsParams []string

	optionalLimitStatsResponse struct {
		required
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
		required{
			CMD: "response",
			ID:  o.ID(),
			To:  request.ID,
		},
		optionalLimitStatsResult{
			Free:  540246016,
			Limit: "10%",
		},
	})
}

type (
	optionalFileInfoRequest struct {
		required
		Params optionalFileInfoParams `json:"params"`
	}
	optionalFileInfoParams struct {
		InnerPath string `json:"inner_path"`
	}

	optionalFileInfoResponse struct {
		required
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
		required{
			CMD: "response",
			ID:  o.ID(),
			To:  request.ID,
		},
		info,
	})
}

type (
	optionalHelpListRequest struct {
		required
		Params map[string]string `json:"params"`
	}

	optionalHelpListResponse struct {
		required
		Result map[string]string `json:"result"`
	}
)

func (o *optionalManager) optionalHelpList(w pluginWriter, site *site.Site, message []byte) error {
	request := new(optionalHelpListRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}

	params := site.Settings.OptionalHelp
	if params == nil {
		params = make(map[string]string)
	}

	return w.WriteJSON(optionalHelpListResponse{
		required{
			CMD: "response",
			ID:  o.ID(),
			To:  request.ID,
		},
		params,
	})
}
