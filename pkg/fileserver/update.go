package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	updateRequest struct {
		CMD    string       `msgpack:"cmd"`
		ReqID  int          `msgpack:"req_id"`
		Params updateParams `msgpack:"params"`
	}
	updateParams struct {
		Site      string `msgpack:"site"`
		InnerPath string `msgpack:"inner_path"`
		Body      []byte `msgpack:"body"`
	}

	updateResponse struct {
		CMD   string `msgpack:"cmd"`
		To    int    `msgpack:"to"`
		Ok    bool   `msgpack:"ok"`
		Error string `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func Update(conn net.Conn, site, innerPath string) (*updateResponse, error) {
	// TODO: include content.json body + diffs
	encoded, err := msgpack.Marshal(&updateRequest{
		CMD:   "update",
		ReqID: 1,
		Params: updateParams{
			Site:      site,
			InnerPath: innerPath,
			Body:      []byte(`{"modified": 0}`),
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(updateResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func updateHandler(conn net.Conn, decoder requestDecoder) error {
	// TODO: validate body and update content
	var r updateRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	data, err := msgpack.Marshal(&updateResponse{
		CMD: "response",
		To:  r.ReqID,
		Ok:  true,
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}