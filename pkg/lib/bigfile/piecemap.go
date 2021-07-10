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

// Implements DecodeMsgpack interface
func (h *hash) DecodeMsgpack(dec *msgpack.Decoder) error {
	var bytes []byte
	if err := dec.Decode(&bytes); err != nil {
		return err
	}
	*h = hash(hex.EncodeToString(bytes))
	return nil
}

// Hashes returns the SHA512/256 hashes for the specified file.
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

// ParsePieceMap parses a messagepack encoded piecemap as
// specified in the documentation bellow:
// https://zeronet.io/docs/help_zeronet/network_protocol/#bigfile-piecemap
func ParsePieceMap(r io.Reader) (*PieceMap, error) {
	pieceMap := new(PieceMap)
	return pieceMap, msgpack.NewDecoder(r).Decode(&pieceMap)
}

// MarshalPieceMap creates a msgpack encoded piecemap.
func MarshalPieceMap(filename string, hashes []string) ([]byte, error) {
	p := make(map[string]map[string][][]byte)
	p[filename] = make(map[string][][]byte)
	for _, hash := range hashes {
		digest, err := hex.DecodeString(hash)
		if err != nil {
			return nil, err
		}
		p[filename]["sha512_pieces"] = append(p[filename]["sha512_pieces"], digest)
	}
	return msgpack.Marshal(p)
}
