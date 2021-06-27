package bigfile

import (
	"encoding/hex"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

type PieceMap map[string]map[string][]hash

type hash []string

var _ msgpack.CustomDecoder = (*hash)(nil)

func (h *hash) DecodeMsgpack(dec *msgpack.Decoder) error {
	var bytes []byte
	if err := dec.Decode(&bytes); err != nil {
		return err
	}
	*h = append(*h, hex.EncodeToString(bytes))
	return nil
}

func ParsePieceMap(r io.Reader) (*PieceMap, error) {
	pieceMap := new(PieceMap)
	return pieceMap, msgpack.NewDecoder(r).Decode(&pieceMap)
}
