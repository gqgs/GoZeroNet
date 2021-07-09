package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	getPieceFieldsRequest struct {
		CMD    string               `msgpack:"cmd"`
		ReqID  int64                `msgpack:"req_id"`
		Params getPieceFieldsParams `msgpack:"params"`
	}
	getPieceFieldsParams struct {
		Site string `msgpack:"site"`
	}

	getPieceFieldsResponse struct {
		CMD               string            `msgpack:"cmd"`
		To                int64             `msgpack:"to"`
		PieceFieldsPacked map[string]string `msgpack:"piecefields_packed"`
	}
)

func (s *server) getPieceFieldsHandler(conn net.Conn, decoder requestDecoder) error {
	var r getPieceFieldsRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	// TODO: implement me

	data, err := msgpack.Marshal(&getPieceFieldsResponse{
		CMD:               "response",
		To:                r.ReqID,
		PieceFieldsPacked: make(map[string]string),
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
