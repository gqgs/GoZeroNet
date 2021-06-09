package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	getFileRequest struct {
		CMD    string        `msgpack:"cmd"`
		ReqID  int           `msgpack:"req_id"`
		Params getFileParams `msgpack:"params"`
	}
	getFileParams struct {
		Site      string `msgpack:"site"`
		InnerPath string `msgpack:"inner_path"`
		Location  int    `msgpack:"location"` // offset location for range requests
		FileSize  int    `msgpack:"file_size,omitempty"`
	}

	getFileResponse struct {
		CMD      string `msgpack:"cmd"`
		To       int    `msgpack:"to"`
		Body     []byte `msgpack:"body"`
		Location int    `msgpack:"location"` // offset location of the last byte sent
		Size     int    `msgpack:"size"`
	}
)

func GetFile(conn net.Conn, site, innerPath string, location, size int) (*getFileResponse, error) {
	encoded, err := msgpack.Marshal(&getFileRequest{
		CMD:   "getFile",
		ReqID: 1,
		Params: getFileParams{
			Site:      site,
			InnerPath: innerPath,
			Location:  location,
			FileSize:  size,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(getFileResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func getFileHandler(conn net.Conn, decoder requestDecoder) error {
	// TODO: get values from storage + handle reputation.
	// Max 512 bytes sent in a request
	var r getFileRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	data, err := msgpack.Marshal(&getFileResponse{
		CMD:      "response",
		To:       r.ReqID,
		Body:     []byte("hello world"),
		Location: 0,
		Size:     123456,
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
