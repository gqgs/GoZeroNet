package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	getHashfieldRequest struct {
		CMD    string             `msgpack:"cmd"`
		ReqID  int64              `msgpack:"req_id"`
		Params getHashfieldParams `msgpack:"params"`
	}
	getHashfieldParams struct {
		Site string `msgpack:"site"`
	}

	getHashfieldResponse struct {
		CMD          string `msgpack:"cmd"`
		To           int64  `msgpack:"to"`
		HashfieldRaw string `msgpack:"hashfield_raw"`
	}
)

func (s *server) getHashfieldHandler(conn net.Conn, decoder requestDecoder) error {
	s.log.Debug("new getHashfield file request")
	var r getHashfieldRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	// TODO: implement me

	data, err := msgpack.Marshal(&getHashfieldResponse{
		CMD:          "response",
		To:           r.ReqID,
		HashfieldRaw: "",
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
