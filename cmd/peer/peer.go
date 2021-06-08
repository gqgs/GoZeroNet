package peer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/file"
)

func ping(addr string) error {
	conn, err := file.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < 5; i++ {
		resp, err := file.Ping(conn)
		if err != nil {
			return err
		}
		jsonDump(resp)
	}

	return nil
}

func handshake(addr string) error {
	rand.Seed(time.Now().UnixNano())

	fileServer, err := file.NewServer(config.RandomIPv4Addr)
	if err != nil {
		return err
	}
	defer fileServer.Shutdown()

	conn, err := file.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := file.Handshake(conn, addr, fileServer)
	jsonDump(resp)
	return err
}

func getFile(addr, site, innerPath string) error {
	conn, err := file.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := file.GetFile(conn, site, innerPath)
	jsonDump(resp)
	return err
}

func jsonDump(v interface{}) {
	d, _ := json.Marshal(v)
	fmt.Println(string(d))
}
