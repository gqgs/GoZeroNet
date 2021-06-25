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
	streamFileRequest struct {
		CMD    string           `msgpack:"cmd"`
		ReqID  int64            `msgpack:"req_id"`
		Params streamFileParams `msgpack:"params"`
	}
	streamFileParams struct {
		Site      string `msgpack:"site"`
		InnerPath string `msgpack:"inner_path"`
		Location  int    `msgpack:"location"` // offset location for range requests
		FileSize  int    `msgpack:"file_size,omitempty"`
	}

	streamFileResponse struct {
		CMD         string `msgpack:"cmd"`
		To          int64  `msgpack:"to"`
		StreamBytes int    `msgpack:"stream_bytes"`
		Error       string `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

// StreamFileFull requests and concatenates all chunks of the file
func StreamFileFull(conn net.Conn, site, innerPath string, size int) (*getFileResponse, error) {
	var body []byte
	var location int

	for {
		_, reader, err := StreamFile(conn, site, innerPath, location, size)
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
		location += len(body)
	}

	return &getFileResponse{
		CMD:      "response",
		Body:     body,
		Location: location,
		Size:     len(body),
	}, nil
}

func StreamFile(conn net.Conn, site, innerPath string, location, size int) (*streamFileResponse, io.Reader, error) {
	encoded, err := msgpack.Marshal(&streamFileRequest{
		CMD:   "streamFile",
		ReqID: counter(),
		Params: streamFileParams{
			Site:      site,
			InnerPath: innerPath,
			Location:  location,
			FileSize:  size,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, nil, err
	}

	// Msgpack's decoder might read more data than it needs to decode
	// the contents of the message. Return the buffered data to include
	// any extraneous content that was read by the decoder.
	result := new(streamFileResponse)
	decoder := msgpack.NewDecoder(conn)
	if err = decoder.Decode(result); err != nil {
		return nil, nil, err
	}
	reader := io.MultiReader(decoder.Buffered(), conn)
	return result, io.LimitReader(reader, int64(result.StreamBytes)), nil
}

func (s *server) streamFileHandler(conn net.Conn, decoder requestDecoder) error {
	s.log.Debug("new streamFile request")
	var r streamFileRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	var size int
	var body []byte
	if info, err := s.contentDB.FileInfo(r.Params.Site, r.Params.InnerPath); err != nil {
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
		size, err = file.ReadAt(body, int64(r.Params.Location))
		if err != nil && err != io.EOF {
			return err
		}
		body = body[:size]

		info.Uploaded += size
		event.BroadcastFileInfoUpdate(r.Params.Site, s.pubsubManager, info)
	}

	data, err := msgpack.Marshal(&streamFileResponse{
		CMD:         "response",
		To:          r.ReqID,
		StreamBytes: size,
	})
	if err != nil {
		return err
	}

	if _, err = conn.Write(data); err != nil {
		return err
	}

	_, err = conn.Write(body)
	return err
}
