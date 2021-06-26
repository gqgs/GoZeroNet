package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	getFileRequest struct {
		CMD    string        `msgpack:"cmd"`
		ReqID  int64         `msgpack:"req_id"`
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
		To       int64  `msgpack:"to"`
		Body     []byte `msgpack:"body"`
		Location int    `msgpack:"location"` // offset location of the last byte sent
		Size     int    `msgpack:"size"`
		Error    string `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

// GetFile requests and concatenates all chunks of the file
func GetFileFull(conn net.Conn, site, innerPath string, size int) (*getFileResponse, error) {
	var body []byte
	var location int

	for {
		resp, err := GetFile(conn, site, innerPath, location, 0)
		if err != nil {
			return nil, err
		}
		if len(resp.Body) == 0 {
			break
		}
		body = append(body, resp.Body...)
		if len(body) >= resp.Size {
			break
		}
		location = resp.Location
	}

	return &getFileResponse{
		CMD:      "response",
		Body:     body,
		Location: location,
		Size:     len(body),
	}, nil
}

// GetFile read a chunk of the file.
// The return is limited to 512KB.
func GetFile(conn net.Conn, site, innerPath string, location, size int) (*getFileResponse, error) {
	encoded, err := msgpack.Marshal(&getFileRequest{
		CMD:   "getFile",
		ReqID: counter(),
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

func (s *server) getFileHandler(conn net.Conn, decoder requestDecoder) error {
	s.log.Debug("new get file request")
	var r getFileRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	body, size, location, err := s.readChunk(r.Params.Site, r.Params.InnerPath, r.Params.Location)
	if err != nil {
		return err
	}

	data, err := msgpack.Marshal(&getFileResponse{
		CMD:      "response",
		To:       r.ReqID,
		Body:     body,
		Location: location,
		Size:     size,
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
