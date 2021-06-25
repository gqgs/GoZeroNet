package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	pingRequest struct {
		CMD    string              `msgpack:"cmd"`
		ReqID  int64               `msgpack:"req_id"`
		Params map[string]struct{} `msgpack:"params"`
	}
	pingResponse struct {
		CMD   string `msgpack:"cmd"`
		To    int64  `msgpack:"to"`
		Body  string `msgpack:"body"`
		Error string `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func Ping(conn net.Conn) (*pingResponse, error) {
	encoded, err := msgpack.Marshal(&pingRequest{
		CMD:    "ping",
		ReqID:  1,
		Params: make(map[string]struct{}),
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(pingResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func (s *server) pingHandler(conn net.Conn, decoder requestDecoder) error {
	s.log.Debug("new ping request")
	var r pingRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}
	data, err := msgpack.Marshal(&pingResponse{
		CMD:  "response",
		To:   r.ReqID,
		Body: "Pong!",
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
