package peer

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
)

func ping(addr string) error {
	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < 5; i++ {
		resp, err := fileserver.Ping(conn)
		if err != nil {
			return err
		}
		jsonDump(resp)
	}

	return nil
}

func handshake(addr string) error {
	rand.Seed(time.Now().UnixNano())

	srv, err := fileserver.NewServer(config.RandomIPv4Addr)
	if err != nil {
		return err
	}
	defer srv.Shutdown()

	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.Handshake(conn, addr, srv)
	jsonDump(resp)
	return err
}

func getFile(addr, site, innerPath string, location, size int) error {
	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.GetFile(conn, site, innerPath, location, size)
	jsonDump(resp)
	return err
}

func streamFile(addr, site, innerPath string, location, size int) error {
	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, stream, err := fileserver.StreamFile(conn, site, innerPath, location, size)
	if err != nil {
		return err
	}
	jsonDump(resp)

	file, err := io.ReadAll(stream)
	jsonDump(struct {
		File []byte
	}{
		File: file,
	})
	return err
}

func checkPort(addr string, port int) error {
	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.CheckPort(conn, port)
	jsonDump(resp)
	return err
}

func pex(addr, site string, need int) error {
	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.Pex(conn, site, need)
	jsonDump(resp)
	return err
}

func jsonDump(v interface{}) {
	d, _ := json.Marshal(v)
	fmt.Println(string(d))
}
