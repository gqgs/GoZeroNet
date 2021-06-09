package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	pexRequest struct {
		CMD    string    `msgpack:"cmd"`
		ReqID  int       `msgpack:"req_id"`
		Params pexParams `msgpack:"params"`
	}
	pexParams struct {
		Site       string   `msgpack:"site"`
		Peers      [][]byte `msgpack:"peers"`
		PeersOnion [][]byte `msgpack:"peers_onion"`
		Need       int      `msgpack:"need"`
	}

	pexResponse struct {
		CMD        string   `msgpack:"cmd"`
		To         int      `msgpack:"to"`
		Peers      [][]byte `msgpack:"peers"`
		PeersOnion [][]byte `msgpack:"peers_onion"`
		Error      string   `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func Pex(conn net.Conn, site string, need int) (*pexResponse, error) {
	// TODO: include peers that the client has
	encoded, err := msgpack.Marshal(&pexRequest{
		CMD:   "pex",
		ReqID: 1,
		Params: pexParams{
			Site:       site,
			Peers:      [][]byte{},
			PeersOnion: [][]byte{},
			Need:       need,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(pexResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func pexHandler(conn net.Conn, decoder requestDecoder) error {
	var r pexRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	// TODO: include peers that the client has
	data, err := msgpack.Marshal(&pexResponse{
		CMD:        "response",
		To:         r.ReqID,
		Peers:      [][]byte{},
		PeersOnion: [][]byte{},
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
