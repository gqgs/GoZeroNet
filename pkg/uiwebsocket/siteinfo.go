package uiwebsocket

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
	"github.com/gqgs/go-zeronet/pkg/site"
)

type (
	siteInfoRequest struct {
		required
		Params json.RawMessage `json:"params"`
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
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	info, err := w.site.Info()
	if err != nil {
		return err
	}

	params := new(siteInfoParams)
	if err := serialize.JSONUnmarshal(payload.Params, params); err != nil {
		if err := serialize.JSONUnmarshal(payload.Params, &params.FileStatus); err != nil {
			return err
		}
	}
	if len(params.FileStatus) > 0 {
		fileInfo, err := w.site.FileInfo(params.FileStatus)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		} else {
			if fileInfo.IsDownloaded {
				info.Event = []interface{}{"file_done", params.FileStatus}
			}
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
