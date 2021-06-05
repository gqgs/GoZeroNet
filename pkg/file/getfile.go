package file

import (
	"log"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	getFileRequest struct {
		CMD    string        `msgpack:"cmd"`
		ReqID  int           `msgpack:"req_id"`
		Params getFileParams `msgpack:"params"`
	}
	getFileParams struct {
		Site      string `msgpack:"site"`
		InnerPath string `msgpack:"inner_path"`
		Location  string `msgpack:"location"`
		FileSize  *int   `msgpack:"file_size"`
	}

	getFileResponse struct {
		CMD      string `msgpack:"cmd"`
		To       int    `msgpack:"to"`
		Body     string `msgpack:"body"`
		Location string `msgpack:"location"`
		Size     int    `msgpack:"size"`
	}
)

func getFileHandler(w http.ResponseWriter, r getFileRequest) {
	data, err := msgpack.Marshal(&getFileResponse{
		CMD:      "response",
		To:       r.ReqID,
		Body:     "TODO",
		Location: "TODO",
		Size:     42,
	})
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
