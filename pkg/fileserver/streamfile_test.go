package fileserver

import (
	"io"
	"testing"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/connection"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/stretchr/testify/assert"
)

func Test_StreamFile(t *testing.T) {
	pubsubManager := pubsub.NewManager()
	contentDB, err := database.NewContentDatabase()
	if err != nil {
		t.Fatal(err)
	}
	defer contentDB.Close()

	tests := []struct {
		name      string
		site      string
		innerPath string
		wantSize  int
	}{
		{
			"small file",
			"site",
			"smallfile",
			1024 * 256,
		},
		{
			"big file",
			"site",
			"bigfile",
			config.FileGetSizeLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo := &event.FileInfo{
				InnerPath: tt.innerPath,
			}

			if err := contentDB.UpdateFile(tt.site, fileInfo); err != nil {
				t.Fatal(err)
			}

			srv, err := NewServer(config.RandomIPv4Addr, contentDB, pubsubManager)
			if err != nil {
				t.Fatal(err)
			}
			defer srv.Shutdown()
			go srv.Listen()

			conn, err := connection.NewConnection(srv.addr)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			resp, reader, err := StreamFile(conn, tt.site, tt.innerPath, 0, 0, 0)
			if err != nil {
				t.Fatal(err)
			}

			body, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "response", resp.CMD)
			assert.Equal(t, tt.wantSize, len(body))
		})
	}
}
