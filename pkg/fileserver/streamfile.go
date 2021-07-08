package fileserver

import (
	"bufio"
	"io"
	"net"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/vmihailenco/msgpack/v5"
)

type (
	streamFileRequest struct {
		CMD    string           `msgpack:"cmd"`
		ReqID  int64            `msgpack:"req_id"`
		Params streamFileParams `msgpack:"params"`
	}
	streamFileParams struct {
		Site      string `msgpack:"site"`
		InnerPath string `msgpack:"inner_path"`
		Location  int    `msgpack:"location"` // offset location for range requests
		BytesRead int    `msgpack:"read_bytes"`
		FileSize  int    `msgpack:"file_size,omitempty"`
	}

	streamFileResponse struct {
		CMD         string `msgpack:"cmd"`
		To          int64  `msgpack:"to"`
		StreamBytes int    `msgpack:"stream_bytes"`
		Location    int    `msgpack:"location"`
		Size        int    `msgpack:"size"`
		Error       string `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

// StreamAtMost requests and concatenates at most `limit` bytes of the file
// starting at the specified location.
func StreamAtMost(conn net.Conn, site, innerPath string, location, limit, size int) (*getFileResponse, error) {
	_, reader, err := StreamFile(conn, site, innerPath, location, size, limit)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &getFileResponse{
		CMD:      "response",
		Body:     body,
		Location: location,
		Size:     len(body),
	}, nil
}

// StreamFileFull requests and concatenates all chunks of the file
func StreamFileFull(conn net.Conn, site, innerPath string, size int) (*getFileResponse, error) {
	var body []byte
	var location int

	for {
		_, reader, err := StreamFile(conn, site, innerPath, location, size, config.FileGetSizeLimit)
		if err != nil {
			return nil, err
		}
		respBody, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		if len(respBody) == 0 {
			break
		}
		body = append(body, respBody...)
		location += len(respBody)
	}

	return &getFileResponse{
		CMD:      "response",
		Body:     body,
		Location: location,
		Size:     len(body),
	}, nil
}

func StreamFile(conn net.Conn, site, innerPath string, location, size, bytesRead int) (*streamFileResponse, io.Reader, error) {
	encoded, err := msgpack.Marshal(&streamFileRequest{
		CMD:   "streamFile",
		ReqID: counter(),
		Params: streamFileParams{
			Site:      site,
			InnerPath: innerPath,
			Location:  location,
			FileSize:  size,
			BytesRead: bytesRead,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, nil, err
	}

	// By default msgpack's decoder might read more data than it needs to decode
	// the contents of the message. An io.ByteScanner is given to the decode to
	// prevent this behavior.
	result := new(streamFileResponse)
	reader := bufio.NewReader(conn)
	decoder := msgpack.NewDecoder(reader)
	if err = decoder.Decode(result); err != nil {
		return nil, nil, err
	}

	return result, io.LimitReader(reader, int64(result.StreamBytes)), nil
}

func (s *server) streamFileHandler(conn net.Conn, decoder requestDecoder) error {
	var r streamFileRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	info, err := s.contentDB.Info(r.Params.Site, r.Params.InnerPath)
	if err != nil {
		return err
	}

	var streamBytes int
	var location int
	var size int

	if info.GetIsDownloaded() {
		s.log.Debugf("streaming file %s at %d/%d", r.Params.InnerPath, r.Params.Location, r.Params.FileSize)

		filePath := path.Join(config.DataDir, r.Params.Site, safe.CleanPath(r.Params.InnerPath))

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = file.Seek(int64(r.Params.Location), 0); err != nil {
			return err
		}

		streamBytes = r.Params.BytesRead
		size = info.GetSize()

		if size-r.Params.Location < streamBytes {
			streamBytes = size - r.Params.Location
		}
		location = r.Params.Location + streamBytes

		info.AddUploaded(streamBytes)
		info.Update(r.Params.Site, s.pubsubManager)

		reader := io.LimitReader(file, int64(streamBytes))
		defer io.Copy(conn, reader)
	}

	data, err := msgpack.Marshal(&streamFileResponse{
		CMD:         "response",
		To:          r.ReqID,
		StreamBytes: streamBytes,
		Location:    location,
		Size:        size,
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
