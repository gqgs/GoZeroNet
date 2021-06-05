package file

import (
	"log"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	pingRequest struct {
		CMD   string `msgpack:"cmd"`
		ReqID int    `msgpack:"req_id"`
	}
	pingResponse struct {
		CMD  string `msgpack:"cmd"`
		To   int    `msgpack:"to"`
		Body string `msgpack:"body"`
	}
)

func pingHandler(w http.ResponseWriter, r pingRequest) {
	data, err := msgpack.Marshal(&pingResponse{
		CMD:  "response",
		To:   r.ReqID,
		Body: "pong",
	})
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
