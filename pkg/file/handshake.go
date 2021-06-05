package file

import (
	"log"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	handshakeRequest struct {
		CMD    string          `msgpack:"cmd"`
		ReqID  int             `msgpack:"req_id"`
		Params handshakeParams `msgpack:"params"`
	}

	handshakeParams struct {
		Crypt          *string  `msgpack:"crypt"`
		CryptSupported []string `msgpack:"crypt_supported"`
		FileserverPort int      `msgpack:"fileserver_port"`
		Onion          string   `msgpack:"onion"`
		Protocol       string   `msgpack:"protocol"`
		PortOpened     bool     `msgpack:"port_opened"`
		PeerID         string   `msgpack:"peer_id"`
		Rev            int      `msgpack:"rev"`
		TargetIP       string   `msgpack:"target_ip"`
		Version        string   `msgpack:"version"`
	}

	handshakeResponse struct {
		CMD            string   `msgpack:"cmd"`
		To             int      `msgpack:"to"`
		Crypt          *string  `msgpack:"crypt"`
		CryptSupported []string `msgpack:"crypt_supported"`
		FileserverPort int      `msgpack:"fileserver_port"`
		Onion          string   `msgpack:"onion"`
		Protocol       string   `msgpack:"protocol"`
		PortOpened     bool     `msgpack:"port_opened"`
		PeerID         string   `msgpack:"peer_id"`
		Rev            int      `msgpack:"rev"`
		TargetIP       string   `msgpack:"target_ip"`
		Version        string   `msgpack:"version"`
	}
)

func handshakeHandler(w http.ResponseWriter, r handshakeRequest) {
	data, err := msgpack.Marshal(&handshakeResponse{
		CMD:            "response",
		To:             r.ReqID,
		Rev:            2092,
		PortOpened:     false,
		FileserverPort: 43111,
		Protocol:       "v2",
	})
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
