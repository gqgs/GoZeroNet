package file

import (
	"net"
	"net/http"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	pingRequest struct {
		CMD    string              `msgpack:"cmd"`
		ReqID  int                 `msgpack:"req_id"`
		Params map[string]struct{} `msgpack:"params"`
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
		CMD:    "ping",
		ReqID:  1,
		Params: make(map[string]struct{}),
	})
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err = conn.Write(data); err != nil {
		return nil, err
	}

	// TODO: is this the best way?
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	result := new(pingResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}
