package plugin

import (
	"errors"
	"path"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/spf13/cast"
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
		Params []interface{} `json:"params"`
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

func (n *bigFilePlugin) bigfileUploadInit(w pluginWriter, s *site.Site, message []byte) error {
	r := new(bigfileUploadInitRequest)
	if err := jsonUnmarshal(message, r); err != nil {
		return err
	}

	if len(r.Params) < 2 {
		return errors.New("bad request")
	}

	innerPath := cast.ToString(r.Params[0])
	size := cast.ToInt(r.Params[1])
	protocol := "xhr"

	if len(r.Params) >= 3 {
		protocol = cast.ToString(r.Params[2])
	}

	nonce := random.HexString(64)
	pieceSize := 1024 * 1024
	innerPath = safe.CleanPath(innerPath)
	relativePath := path.Join(config.DataDir, s.Address(), innerPath)

	upload := site.Upload{
		Added:     time.Now(),
		Site:      s.Address(),
		InnerPath: innerPath,
		Size:      size,
		PieceSize: pieceSize,
		Piecemap:  innerPath + ".piecemap.msgpack",
	}

	s.UploadInit(upload, nonce)

	var url string
	switch protocol {
	case "xhr":
		url = "/ZeroNet-Internal/BigfileUpload?upload_nonce=" + nonce
	case "websocket":
		url = "{origin}/ZeroNet-Internal/BigfileUploadWebsocket?upload_nonce=" + nonce
	default:
		return errors.New("unknown protocol")
	}

	return w.WriteJSON(bigfileUploadInitResponse{
		required{
			CMD: "response",
			ID:  n.ID(),
			To:  r.ID,
		},
		bigfileUploadInitResult{
			URL:              url,
			PieceSize:        pieceSize,
			InnerPath:        innerPath,
			FileRelativePath: relativePath,
		},
	})
}
