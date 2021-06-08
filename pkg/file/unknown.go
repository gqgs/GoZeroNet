package file

import (
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

type unknownRequest struct {
	ReqID int `msgpack:"req_id"`
}

type unknownResponse struct {
	CMD   string `msgpack:"cmd"`
	To    int    `msgpack:"to"`
	Error string `msgpack:"error"`
}

func unknownHandler(w io.Writer, r unknownRequest) error {
	data, err := msgpack.Marshal(&unknownResponse{
		CMD:   "response",
		To:    r.ReqID,
		Error: "Unknown cmd",
	})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}
