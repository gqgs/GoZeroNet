package random

import (
	"math/rand"
	"strings"
)

type alphabet string

const base62Alphabet alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const hexAlphabet alphabet = "abcdef0123456789"

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
		builder.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return builder.String()
}
