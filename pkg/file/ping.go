package file

import (
	"bytes"
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

func (s *server) pingHandler(w http.ResponseWriter, r pingRequest) {
	data, err := msgpack.Marshal(&pingResponse{
		CMD:  "response",
		To:   r.ReqID,
		Body: "Pong!",
	})
	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (s *server) Ping(addr string) (*pingResponse, error) {
	data, err := msgpack.Marshal(&pingRequest{
		CMD:   "ping",
		ReqID: 1,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, "http://"+addr, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := new(pingResponse)
	return result, msgpack.NewDecoder(resp.Body).Decode(result)
}
