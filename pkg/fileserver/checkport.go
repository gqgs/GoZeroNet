package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	checkPortRequest struct {
		CMD    string          `msgpack:"cmd"`
		ReqID  int             `msgpack:"req_id"`
		Params checkPortParams `msgpack:"params"`
	}
	checkPortParams struct {
		Port int `msgpack:"port"`
	}

	checkPortResponse struct {
		CMD        string `msgpack:"cmd"`
		To         int    `msgpack:"to"`
		Status     string `msgpack:"status"`
		IPExternal string `msgpack:"ip_external"`
		Error      string `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func CheckPort(conn net.Conn, port int) (*checkPortResponse, error) {
	encoded, err := msgpack.Marshal(&checkPortRequest{
		CMD:   "checkport",
		ReqID: 1,
		Params: checkPortParams{
			Port: port,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(checkPortResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func checkPortHandler(conn net.Conn, decoder requestDecoder, server *server) error {
	var r checkPortRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	status := "closed"
	if r.Params.Port == server.port {
		status = "open"
	}
	data, err := msgpack.Marshal(&checkPortResponse{
		CMD:        "response",
		To:         r.ReqID,
		Status:     status,
		IPExternal: conn.RemoteAddr().String(),
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
