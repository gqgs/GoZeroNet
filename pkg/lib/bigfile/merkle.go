package bigfile

import (
	"crypto/sha512"
	"encoding/hex"
)

type hashFunc = func(b []byte) [32]byte

// MerkleRoot calculates the root of the
// merkle tree given the input hashes.
// It considers the tree to the generated according to the
// documentation bellow:
// https://github.com/Tierion/merkle-tools#notes
func MerkleRoot(hashes []string, hasher hashFunc) string {
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
		levels = append(levels, hex.EncodeToString(level[:]))
		hashes = hashes[2:]
	}
	levels = append(levels, hashes...)
	return MerkleRoot(levels, hasher)
}

// TODO: simplify in 1.17
// https://github.com/golang/go/issues/46505
func sha512_256(b []byte) [32]byte {
	sum := sha512.Sum512(b)
	var digest [32]byte
	copy(digest[:], sum[:])
	return digest
}
