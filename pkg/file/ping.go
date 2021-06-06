package file

import (
	"io"
	"net"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
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

func (s *server) pingHandler(w io.Writer, r pingRequest) error {
	data, err := msgpack.Marshal(&pingResponse{
		CMD:  "response",
		To:   r.ReqID,
		Body: "Pong!",
	})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
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
	conn.SetReadDeadline(time.Now().Add(config.ReadDeadline))

	result := new(pingResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}
