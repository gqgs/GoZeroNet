package bigfile

import (
	"encoding/hex"

	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
)

type hashFunc = func(b []byte) string

// MerkleRoot calculates the root of the
// merkle tree given the input hashes.
// It considers the tree to be generated according to the
// documentation bellow:
// https://github.com/Tierion/merkle-tools#notes
func MerkleRoot(hashes []string) string {
	return merkleRoot(hashes, crypto.Sha512_256)
}

func merkleRoot(hashes []string, hasher hashFunc) string {
	if len(hashes) == 0 {
		return ""
	}

	if len(hashes) == 1 {
		return hashes[0]
	}

	var levels []string
	for len(hashes) >= 2 {
		left, _ := hex.DecodeString(hashes[0])
		right, _ := hex.DecodeString(hashes[1])
		level := hasher(append(left, right...))
		levels = append(levels, level)
		hashes = hashes[2:]
	}
	levels = append(levels, hashes...)
	return merkleRoot(levels, hasher)
}
