package fileserver

import (
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	listModifiedRequest struct {
		CMD    string             `msgpack:"cmd"`
		ReqID  int64              `msgpack:"req_id"`
		Params listModifiedParams `msgpack:"params"`
	}
	listModifiedParams struct {
		Site  string `msgpack:"site"`
		Since int    `msgpack:"since"`
	}

	listModifiedResponse struct {
		CMD           string         `msgpack:"cmd"`
		To            int64          `msgpack:"to"`
		ModifiedFiles map[string]int `msgpack:"modified_files"`
		Error         string         `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func ListModified(conn net.Conn, site string, since int) (*listModifiedResponse, error) {
	encoded, err := msgpack.Marshal(&listModifiedRequest{
		CMD:   "listModified",
		ReqID: counter(),
		Params: listModifiedParams{
			Site:  site,
			Since: since,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(listModifiedResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func (s *server) listModifiedHandler(conn net.Conn, decoder requestDecoder) error {
	s.log.Debug("new list modified request")
	var r listModifiedRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	modified, err := s.contentDB.UpdatedContent(r.Params.Site, r.Params.Since)
	if err != nil {
		return err
	}

	data, err := msgpack.Marshal(&listModifiedResponse{
		CMD:           "response",
		To:            r.ReqID,
		ModifiedFiles: modified,
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
