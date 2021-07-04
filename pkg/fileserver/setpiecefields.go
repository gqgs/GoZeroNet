package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	setPieceFieldsRequest struct {
		CMD    string               `msgpack:"cmd"`
		ReqID  int64                `msgpack:"req_id"`
		Params setPieceFieldsParams `msgpack:"params"`
	}
	setPieceFieldsParams struct {
		Site              string            `msgpack:"site"`
		PieceFieldsPacked map[string]string `msgpack:"piecefields_packed"`
	}

	setPieceFieldsResponse struct {
		CMD string `msgpack:"cmd"`
		To  int64  `msgpack:"to"`
		Ok  string `msgpack:"ok"`
	}
)

func (s *server) setPieceFieldsHandler(conn net.Conn, decoder requestDecoder) error {
	s.log.Debug("new setPieceFields file request")
	var r setPieceFieldsRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	// TODO: implement me

	data, err := msgpack.Marshal(&setPieceFieldsResponse{
		CMD: "response",
		To:  r.ReqID,
		Ok:  "Updated",
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
