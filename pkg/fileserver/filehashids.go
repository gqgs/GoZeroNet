package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	findHashIDsRequest struct {
		CMD    string            `msgpack:"cmd"`
		ReqID  int               `msgpack:"req_id"`
		Params findHashIDsParams `msgpack:"params"`
	}
	findHashIDsParams struct {
		Site    string `msgpack:"site"`
		HashIDs []int  `msgpack:"hash_ids"`
	}

	findHashIDsResponse struct {
		CMD        string         `msgpack:"cmd"`
		To         int            `msgpack:"to"`
		Peers      map[int][]byte `msgpack:"peers"`
		PeersOnion map[int][]byte `msgpack:"peers_onion"`
		Error      string         `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func FindHashIDs(conn net.Conn, site string, hashIDs ...int) (*findHashIDsResponse, error) {
	encoded, err := msgpack.Marshal(&findHashIDsRequest{
		CMD:   "findHashIds",
		ReqID: 1,
		Params: findHashIDsParams{
			Site:    site,
			HashIDs: hashIDs,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(findHashIDsResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func findHashIDsHandler(conn net.Conn, decoder requestDecoder) error {
	var r findHashIDsRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}
	data, err := msgpack.Marshal(&findHashIDsResponse{
		CMD:        "response",
		To:         r.ReqID,
		Peers:      make(map[int][]byte),
		PeersOnion: make(map[int][]byte),
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
