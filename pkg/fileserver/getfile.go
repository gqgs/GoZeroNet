package fileserver

import (
	"errors"
	"io"
	"net"
	"os"
	"path"
	"strings"

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

func (s *server) readChunk(site, innerPath string, location int) (body []byte, size, newLocation int, err error) {
	if strings.HasSuffix(innerPath, "content.json") {
		info, err := s.contentDB.ContentInfo(site, innerPath)
		if err != nil {
			if errors.Is(err, database.ErrFileNotFound) {
				return nil, 0, 0, nil
			}
			return nil, 0, 0, err
		}
		body, err = readChunk(site, innerPath, location)
		if err != nil {
			return nil, 0, 0, err
		}
		return body, info.Size, location + len(body), nil
	}

	info, err := s.contentDB.FileInfo(site, innerPath)
	if err != nil {
		if errors.Is(err, database.ErrFileNotFound) {
			return nil, 0, 0, nil
		}
		return nil, 0, 0, err
	}
	if !info.IsDownloaded {
		return nil, 0, 0, nil
	}
	body, err = readChunk(site, innerPath, location)
	if err != nil {
		return nil, 0, 0, err
	}
	info.Uploaded += len(body)
	event.BroadcastFileInfoUpdate(site, s.pubsubManager, info)
	return body, info.Size, location + len(body), nil
}

func readChunk(site, innerPath string, location int) (body []byte, err error) {
	innerPath = path.Join(config.DataDir, site, safe.CleanPath(innerPath))
	file, err := os.Open(innerPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body = make([]byte, config.FileGetSizeLimit)
	read, err := file.ReadAt(body, int64(location))
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body[:read], nil
}
