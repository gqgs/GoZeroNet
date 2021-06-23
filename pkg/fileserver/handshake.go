package fileserver

import (
	"fmt"
	"net"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/vmihailenco/msgpack/v5"
)

type (
	handshakeRequest struct {
		CMD    string          `msgpack:"cmd"`
		ReqID  int             `msgpack:"req_id"`
		Params handshakeParams `msgpack:"params"`
	}

	handshakeParams struct {
		Crypt          string   `msgpack:"crypt"`
		CryptSupported []string `msgpack:"crypt_supported"`
		FileserverPort int      `msgpack:"fileserver_port"`
		Onion          string   `msgpack:"onion"`
		Protocol       string   `msgpack:"protocol"`
		PortOpened     bool     `msgpack:"port_opened"`
		PeerID         string   `msgpack:"peer_id"`
		Rev            int      `msgpack:"rev"`
		TargetIP       string   `msgpack:"target_ip"`
		UseBinType     bool     `msgpack:"use_bin_type"`
		Version        string   `msgpack:"version"`
	}

	handshakeResponse struct {
		CMD            string   `msgpack:"cmd"`
		To             int      `msgpack:"to"`
		Crypt          string   `msgpack:"crypt"`
		CryptSupported []string `msgpack:"crypt_supported"`
		FileserverPort int      `msgpack:"fileserver_port"`
		Onion          string   `msgpack:"onion"`
		Protocol       string   `msgpack:"protocol"`
		PortOpened     bool     `msgpack:"port_opened"`
		PeerID         string   `msgpack:"peer_id"`
		Rev            int      `msgpack:"rev"`
		TargetIP       string   `msgpack:"target_ip"`
		UseBinType     bool     `msgpack:"use_bin_type"`
		Version        string   `msgpack:"version"`
		Error          string   `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func Handshake(conn net.Conn, addr string) (*handshakeResponse, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, addr)
	}

	var peerID string
	if peer, ok := conn.(interface {
		ID() string
	}); ok {
		peerID = peer.ID()
	}

	encoded, err := msgpack.Marshal(&handshakeRequest{
		CMD:   "handshake",
		ReqID: 1,
		Params: handshakeParams{
			CryptSupported: make([]string, 0),
			FileserverPort: config.FileServerPort,
			Protocol:       config.Protocol,
			PortOpened:     config.PortOpened,
			PeerID:         peerID,
			Rev:            config.Rev,
			UseBinType:     config.UseBinType,
			Version:        config.Version,
			TargetIP:       host,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(handshakeResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func handshakeHandler(conn net.Conn, decoder requestDecoder, fileServer *server) error {
	var r handshakeRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	var peerID string
	if peer, ok := conn.(interface {
		ID() string
	}); ok {
		peerID = peer.ID()
	}

	encoded, err := msgpack.Marshal(&handshakeResponse{
		CMD:            "response",
		To:             r.ReqID,
		CryptSupported: make([]string, 0),
		FileserverPort: fileServer.port,
		Protocol:       config.Protocol,
		PortOpened:     config.PortOpened,
		PeerID:         peerID,
		Rev:            config.Rev,
		UseBinType:     config.UseBinType,
		Version:        config.Version,
		TargetIP:       host,
	})
	if err != nil {
		return err
	}
	_, err = conn.Write(encoded)
	return err
}
