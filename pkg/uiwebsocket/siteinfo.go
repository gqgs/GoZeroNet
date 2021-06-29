package uiwebsocket

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	siteInfoRequest struct {
		required
		Params siteInfoParams `json:"params"`
	}

	siteInfoParams struct {
		FileStatus string `json:"file_status"`
	}
	siteInfoResponse struct {
		required
		Result *site.Info `json:"result"`
	}
)

func (w *uiWebsocket) siteInfo(rawMessage []byte, message Message) error {
	payload := new(siteInfoRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
		return err
	}

	info, err := w.site.Info()
	if err != nil {
		return err
	}

	fileInfo, err := w.site.FileInfo(payload.Params.FileStatus)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		if fileInfo.IsDownloaded {
			info.Event = []interface{}{"file_done", payload.Params.FileStatus}
		}
	}

	return w.conn.WriteJSON(siteInfoResponse{
		required{
			CMD: "response",
			To:  message.ID,
			ID:  w.ID(),
		},
		info,
	})
}
