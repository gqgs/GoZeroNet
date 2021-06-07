package file

import (
	"io"
	"net"
	"time"

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
	}
)

func (s *server) Handshake(addr string) (*handshakeResponse, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	encoded, err := msgpack.Marshal(&handshakeRequest{
		CMD:   "handshake",
		ReqID: 1,
		Params: handshakeParams{
			Crypt:          "tls-rsa",
			CryptSupported: []string{"tls-rsa"},
			FileserverPort: s.port,
			Protocol:       config.Protocol,
			PortOpened:     config.PortOpened,
			PeerID:         s.peerID,
			Rev:            config.Rev,
			UseBinType:     config.UseBinType,
			Version:        config.Version,
			TargetIP:       host,
		},
	})
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(config.Deadline))

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(handshakeResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)

}

func (s *server) handshakeHandler(w io.Writer, r handshakeRequest) error {
	// TODO: This will panic if writer doesn't implement net.Conn
	// find a better way to get the remote host here
	host, _, err := net.SplitHostPort(w.(net.Conn).RemoteAddr().String())
	if err != nil {
		return err
	}

	encoded, err := msgpack.Marshal(&handshakeResponse{
		CMD:            "response",
		To:             r.ReqID,
		Crypt:          "tls-rsa",
		CryptSupported: []string{"tls-rsa"},
		FileserverPort: s.port,
		Protocol:       config.Protocol,
		PortOpened:     config.PortOpened,
		PeerID:         s.peerID,
		Rev:            config.Rev,
		UseBinType:     config.UseBinType,
		Version:        config.Version,
		TargetIP:       host,
	})
	if err != nil {
		return err
	}
	_, err = w.Write(encoded)
	return err
}
