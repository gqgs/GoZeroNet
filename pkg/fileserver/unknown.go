package fileserver

import (
	"net"

	"github.com/spf13/cast"
	"github.com/vmihailenco/msgpack/v5"
)

type unknownResponse struct {
	CMD   string `msgpack:"cmd"`
	To    int    `msgpack:"to"`
	Error string `msgpack:"error"`
}

func (s *server) unknownHandler(conn net.Conn, decoder requestDecoder) error {
	s.log.Debug("new unknown request")
	reqID, err := decodeKey(decoder, "req_id")
	if err != nil {
		return err
	}

	data, err := msgpack.Marshal(&unknownResponse{
		CMD:   "response",
		To:    cast.ToInt(reqID),
		Error: "Unknown cmd",
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
