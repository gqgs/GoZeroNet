package file

import (
	"log"
	"net/http"

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

func (s *server) handshakeHandler(w http.ResponseWriter, r handshakeRequest) {
	data, err := msgpack.Marshal(&handshakeResponse{
		CMD:            "response",
		To:             r.ReqID,
		Crypt:          "tls-rsa",
		CryptSupported: []string{"tls-rsa"},
		FileserverPort: config.FileServer.Port,
		Protocol:       config.Protocol,
		PortOpened:     config.PortOpened,
		PeerID:         s.peerID,
		Rev:            config.Rev,
		UseBinType:     config.UseBinType,
		Version:        config.Version,
	})
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
