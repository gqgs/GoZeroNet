package fileserver

import (
	"encoding/binary"
	"net"

	"github.com/gqgs/go-zeronet/pkg/lib/ip"
	"github.com/vmihailenco/msgpack/v5"
)

type (
	pexRequest struct {
		CMD    string    `msgpack:"cmd"`
		ReqID  int64     `msgpack:"req_id"`
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
		To         int64    `msgpack:"to"`
		Peers      [][]byte `msgpack:"peers"`
		PeersOnion [][]byte `msgpack:"peers_onion"`
		Error      string   `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func Pex(conn net.Conn, site string, need int) (*pexResponse, error) {
	encoded, err := msgpack.Marshal(&pexRequest{
		CMD:   "pex",
		ReqID: counter(),
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

func (s *server) pexHandler(conn net.Conn, decoder requestDecoder) error {
	var r pexRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	var peers [][]byte
	peerList, err := s.contentDB.Peers(r.Params.Site, r.Params.Need)
	if err != nil {
		return err
	}
	for _, peer := range peerList {
		packed := ip.PackIPv4(peer, binary.LittleEndian)
		if packed != nil {
			peers = append(peers, packed)
		}
	}

	data, err := msgpack.Marshal(&pexResponse{
		CMD:        "response",
		To:         r.ReqID,
		Peers:      peers,
		PeersOnion: [][]byte{},
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
