package fileserver

import (
	"fmt"
	"net"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/vmihailenco/msgpack/v5"
)

type (
	updateRequest struct {
		CMD    string       `msgpack:"cmd"`
		ReqID  int64        `msgpack:"req_id"`
		Params updateParams `msgpack:"params"`
	}
	updateParams struct {
		Site      string `msgpack:"site"`
		InnerPath string `msgpack:"inner_path"`
		Body      []byte `msgpack:"body"`
	}

	updateResponse struct {
		CMD   string `msgpack:"cmd"`
		To    int64  `msgpack:"to"`
		Ok    string `msgpack:"ok"`
		Error string `msgpack:"error,omitempty" json:"error,omitempty"`
	}
)

func Update(conn net.Conn, site, innerPath string) (*updateResponse, error) {
	filePath := path.Join(config.DataDir, site, safe.CleanPath(innerPath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	encoded, err := msgpack.Marshal(&updateRequest{
		CMD:   "update",
		ReqID: counter(),
		Params: updateParams{
			Site:      site,
			InnerPath: innerPath,
			Body:      content,
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(encoded); err != nil {
		return nil, err
	}

	result := new(updateResponse)
	return result, msgpack.NewDecoder(conn).Decode(result)
}

func (s *server) updateHandler(conn net.Conn, decoder requestDecoder) error {
	var r updateRequest
	if err := decoder.Decode(&r); err != nil {
		return err
	}

	event.BroadcastSiteUpdate(r.Params.Site, s.pubsubManager, &event.SiteUpdate{
		InnerPath: r.Params.InnerPath,
		Body:      r.Params.Body,
	})

	data, err := msgpack.Marshal(&updateResponse{
		CMD: "response",
		To:  r.ReqID,
		Ok:  fmt.Sprintf("Thanks, file %s updated!", r.Params.InnerPath),
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
