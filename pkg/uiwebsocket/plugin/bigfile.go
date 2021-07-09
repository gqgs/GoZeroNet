package plugin

import (
	"encoding/json"

	"github.com/gqgs/go-zeronet/pkg/site"
)

type bigFilePlugin struct {
	ID IDFunc
}

func NewBigFile(idFunc IDFunc) Plugin {
	return &bigFilePlugin{
		ID: idFunc,
	}
}

func (*bigFilePlugin) Name() string {
	return "Bigfile"
}

func (*bigFilePlugin) Description() string {
	return "Manage big files"
}

func (n *bigFilePlugin) Handler(cmd string) (HandlerFunc, bool) {
	switch cmd {
	case "bigfileUploadInit":
		return n.bigfileUploadInit, true
	default:
		return nil, false
	}
}

type (
	bigfileUploadInitRequest struct {
		required
		Params bigfileUploadInitParams `json:"params"`
	}
	bigfileUploadInitParams struct {
		InnerPath string `json:"inner_path"`
		Size      int    `json:"size"`
	}

	bigfileUploadInitResponse struct {
		required
		Result bigfileUploadInitResult `json:"result"`
	}

	bigfileUploadInitResult struct {
		URL              string `json:"url"`
		PieceSize        int    `json:"piece_size"`
		InnerPath        string `json:"inner_path"`
		FileRelativePath string `json:"file_relative_path"`
	}
)

func (n *bigFilePlugin) bigfileUploadInit(w pluginWriter, site *site.Site, message []byte) error {
	request := new(bigfileUploadInitRequest)
	if err := json.Unmarshal(message, request); err != nil {
		return err
	}
	return w.WriteJSON(bigfileUploadInitResponse{
		required{
			CMD: "response",
			ID:  n.ID(),
			To:  request.ID,
		},
		bigfileUploadInitResult{},
	})
}
