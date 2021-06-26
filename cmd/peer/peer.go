package peer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
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
	srv, err := fileserver.NewServer(config.RandomIPv4Addr, nil, nil)
	if err != nil {
		return err
	}
	defer srv.Shutdown()

	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.Handshake(conn, addr)
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

func listModified(addr, site string, since int) error {
	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.ListModified(conn, site, int64(since))
	jsonDump(resp)
	return err
}

func update(addr, site, innerPath string) error {
	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.Update(conn, site, innerPath)
	jsonDump(resp)
	return err
}

func jsonDump(v interface{}) {
	d, _ := json.Marshal(v)
	fmt.Println(string(d))
}

func findHashIDs(addr, site string, hashList ...string) error {
	ids := make([]int, len(hashList))
	for i, hash := range hashList {
		id, err := crypto.HashID(hash)
		if err != nil {
			return err
		}
		ids[i] = id
	}

	conn, err := connection.NewConnection(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := fileserver.FindHashIDs(conn, site, ids...)
	jsonDump(resp)
	return err
}
