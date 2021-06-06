package peer

import (
	"context"
	"fmt"

	"github.com/gqgs/go-zeronet/pkg/file"
)

func ping(ctx context.Context, peer string) error {
	fileServer := file.NewServer()
	resp, err := fileServer.Ping(peer)
	fmt.Printf("%+v\n", resp)
	return err
}

func getFile(ctx context.Context) error {
	panic("implement me")
}

func cmd(ctx context.Context) error {
	panic("implement me")
}
