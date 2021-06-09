package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	listModifiedRequest struct {
		CMD    string             `msgpack:"cmd"`
		ReqID  int                `msgpack:"req_id"`
		Params listModifiedParams `msgpack:"params"`
	}
	listModifiedParams struct {
		Site  string `msgpack:"site"`
		Since int    `msgpack:"since"`
	}

	listModifiedResponse struct {
		CMD           string         `msgpack:"cmd"`
		To            int            `msgpack:"to"`
		ModifiedFiles map[string]int `msgpack:"modified_files"`
	}
)

func ListModified(conn net.Conn, site string, since int) (*listModifiedResponse, error) {
	encoded, err := msgpack.Marshal(&listModifiedRequest{
		CMD:   "listModified",
		ReqID: 1,
		Params: listModifiedParams{
			Site:  site,
			Since: since,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(listModifiedResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func listModifiedHandler(conn net.Conn, decoder requestDecoder) error {
	var r listModifiedRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	// TODO: list modified files
	data, err := msgpack.Marshal(&listModifiedResponse{
		CMD:           "response",
		To:            r.ReqID,
		ModifiedFiles: make(map[string]int),
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
