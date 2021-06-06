package random

import (
	"crypto/rand"
	"math/big"
	"strings"
)

type alphabet string

const base62Alphabet alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const hexAlphabet alphabet = "abcdef0123456789"

// Returns a random byte contained in the alphabet.
// This method panics if there is a problem generating an int
// from reading crypto/rand.Reader.
func (a alphabet) randomByte() byte {
	max := big.NewInt(int64(len(a)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return a[n.Int64()]
}

// Base62String returns a random base62 encoded string of size `length`
func Base62String(length int) string {
	return generateString(length, base62Alphabet)
}

// HexString returns a random hex encoded string of size `length`
func HexString(length int) string {
	return generateString(length, hexAlphabet)
}

func generateString(length int, alphabet alphabet) string {
	builder := new(strings.Builder)
	builder.Grow(length)
	for i := 0; i < length; i++ {
		builder.WriteByte(alphabet.randomByte())
	}
	return builder.String()
}

// PeerID returns an ID that can be used to identify connections.
// "Setting the Peer ID to "UT3530" tells trackers that you're using uTorrent v3.5.3"
// https://github.com/jaruba/PowderWeb/wiki/Guide#private-torrent-trackers
func PeerID() string {
	return "-UT3530-" + Base62String(12)
}
