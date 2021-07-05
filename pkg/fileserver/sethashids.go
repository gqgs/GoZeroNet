package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	setHashfieldRequest struct {
		CMD    string             `msgpack:"cmd"`
		ReqID  int64              `msgpack:"req_id"`
		Params setHashfieldParams `msgpack:"params"`
	}
	setHashfieldParams struct {
		Site         string `msgpack:"site"`
		HashfieldRaw []byte `msgpack:"hashfield_raw"`
	}

	setHashfieldResponse struct {
		CMD string `msgpack:"cmd"`
		To  int64  `msgpack:"to"`
		Ok  string `msgpack:"ok"`
	}
)

func (s *server) setHashfieldHandler(conn net.Conn, decoder requestDecoder) error {
	var r setHashfieldRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	// TODO: implement me

	data, err := msgpack.Marshal(&setHashfieldResponse{
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
