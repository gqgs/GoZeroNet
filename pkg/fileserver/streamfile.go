package fileserver

import (
	"io"
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	streamFileRequest struct {
		CMD    string           `msgpack:"cmd"`
		ReqID  int              `msgpack:"req_id"`
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
		To          int    `msgpack:"to"`
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
		To:       1,
		Body:     body,
		Location: location,
		Size:     len(body),
	}, nil
}

func StreamFile(conn net.Conn, site, innerPath string, location, size int) (*streamFileResponse, io.Reader, error) {
	encoded, err := msgpack.Marshal(&streamFileRequest{
		CMD:   "streamFile",
		ReqID: 1,
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

func streamFileHandler(conn net.Conn, decoder requestDecoder) error {
	// TODO: get values from storage + handle reputation + write to stream
	var r streamFileRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	file := []byte("hello wolrd")

	encoded, err := msgpack.Marshal(&streamFileResponse{
		CMD:         "response",
		To:          r.ReqID,
		StreamBytes: len(file),
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(append(encoded, file...))
	return err
}
