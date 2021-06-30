package bigfile

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePieceMap(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     *PieceMap
	}{
		{
			"image piecemap",
			"testdata/image.jpg.piecemap.msgpack",
			&PieceMap{
				"3ff9dd82ccd831865973a943642ff07ee32d8344.jpg": {
					"sha512_pieces": {
						hash("9925a54fe7fe03488e4bbdeddef906f9353b763f7a1f483653360901b6c7e5bb"),
						hash("a148621b234397a6b347f1e84c3b6094e13094e356942e87a1dc7cc647595216"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := os.Open(tt.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()

			pieceMap, err := ParsePieceMap(reader)
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(pieceMap, tt.want))
		})
	}
}
