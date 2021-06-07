package peer

import (
	"context"
	"encoding/json"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/file"
)

func ping(ctx context.Context, addr string) error {
	fileServer, err := file.NewServer(config.FileServerAddr)
	if err != nil {
		return err
	}
	resp, err := fileServer.Ping(addr)
	dump(resp)
	return err
}

func handshake(ctx context.Context, addr string) error {
	fileServer, err := file.NewServer(config.FileServerAddr)
	if err != nil {
		return err
	}
	resp, err := fileServer.Handshake(addr)
	dump(resp)
	return err
}

func getFile(ctx context.Context, addr, site, innerPath string) error {
	fileServer, err := file.NewServer(config.FileServerAddr)
	if err != nil {
		return err
	}
	resp, err := fileServer.GetFile(addr, site, innerPath)
	dump(resp)
	return err
}

// Dumps v in a easy to read format
func dump(v interface{}) {
	d, _ := json.MarshalIndent(v, "", " ")
	println(string(d))
}
