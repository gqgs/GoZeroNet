package peer

import (
	"context"
	"fmt"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/file"
)

func ping(ctx context.Context, addr string) error {
	fileServer, err := file.NewServer(config.FileServerAddr)
	if err != nil {
		return err
	}
	resp, err := fileServer.Ping(addr)
	fmt.Printf("%+v\n", resp)
	return err
}

func handshake(ctx context.Context, addr string) error {
	fileServer, err := file.NewServer(config.FileServerAddr)
	if err != nil {
		return err
	}
	resp, err := fileServer.Handshake(addr)
	fmt.Printf("%+v\n", resp)
	return err
}

func getFile(ctx context.Context, addr, site, innerPath string) error {
	fileServer, err := file.NewServer(config.FileServerAddr)
	if err != nil {
		return err
	}
	resp, err := fileServer.GetFile(addr, site, innerPath)
	fmt.Printf("%+v\n", resp)
	return err
}
