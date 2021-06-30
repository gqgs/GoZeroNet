package bigfile

import (
	"encoding/hex"
	"errors"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

type PieceMap map[string]map[string][]hash

type hash string

var _ msgpack.CustomDecoder = (*hash)(nil)

func (h *hash) DecodeMsgpack(dec *msgpack.Decoder) error {
	var bytes []byte
	if err := dec.Decode(&bytes); err != nil {
		return err
	}
	*h = hash(hex.EncodeToString(bytes))
	return nil
}

func (s *PieceMap) Hashes(file string) ([]string, error) {
	var hashes []string

	files, ok := (*s)[file]
	if !ok {
		return nil, errors.New("file not found in piecemap")
	}
	for _, h := range files["sha512_pieces"] {
		hashes = append(hashes, string(h))
	}

	return hashes, nil
}

func ParsePieceMap(r io.Reader) (*PieceMap, error) {
	pieceMap := new(PieceMap)
	return pieceMap, msgpack.NewDecoder(r).Decode(&pieceMap)
}
