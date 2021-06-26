package fileserver

import (
	"errors"
	"io"
	"net"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
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

	var size int
	var body []byte
	var location int

	info, err := s.contentDB.FileInfo(r.Params.Site, r.Params.InnerPath)
	if err != nil {
		if !errors.Is(err, database.ErrFileNotFound) {
			return err
		}
	} else {
		innerPath := path.Join(config.DataDir, r.Params.Site, safe.CleanPath(r.Params.InnerPath))
		file, err := os.Open(innerPath)
		if err != nil {
			return err
		}
		defer file.Close()

		body = make([]byte, config.FileGetSizeLimit)
		read, err := file.ReadAt(body, int64(r.Params.Location))
		if err != nil && err != io.EOF {
			return err
		}
		body = body[:read]
		size = info.Size
		location = r.Params.Location + read
		info.Uploaded += read
		event.BroadcastFileInfoUpdate(r.Params.Site, s.pubsubManager, info)
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
