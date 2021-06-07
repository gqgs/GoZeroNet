package peer

import (
	"encoding/json"
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
		dump(resp)
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
	dump(resp)
	return err
}

func getFile(addr, site, innerPath string) error {
	conn, err := file.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := file.GetFile(conn, site, innerPath)
	dump(resp)
	return err
}

// Dumps v in a easy to read format
func dump(v interface{}) {
	d, _ := json.Marshal(v)
	println(string(d))
}
