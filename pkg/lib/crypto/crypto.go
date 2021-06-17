package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/ripemd160"
)

// IsValidSignature return true if message was signed with key related to address
func IsValidSignature(message []byte, base64Sign, address string) bool {
	pubkey, err := RecoverPublicKey(message, base64Sign)
	if err != nil {
		return false
	}
	return address == PublicKeyToAddress(pubkey)
}

// RecoverPublicKey recovers the public key given a message and the message signature
func RecoverPublicKey(message []byte, base64Sign string) ([]byte, error) {
	sign, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return nil, err
	}

	sign, err = coincurveSig(sign)
	if err != nil {
		return nil, err
	}

	return secp256k1.RecoverPubkey(hash(message), sign)
}

// PublicKeyToAddress converts a public key to ZN style addresses
func PublicKeyToAddress(pubKey []byte) string {
	sha256Digest := sha256.Sum256(pubKey)
	ripemdHasher := ripemd160.New()
	if _, err := ripemdHasher.Write(sha256Digest[:]); err != nil {
		return ""
	}
	digest := ripemdHasher.Sum(nil)
	digest = append([]byte{0}, digest...)
	result := base58CheckEncode(digest)
	return result
}

// base58 has a CheckEncode method but it prepends a version
// which is a different behavior from the Python lib
func base58CheckEncode(input []byte) string {
	b := make([]byte, 0, len(input)+4)
	b = append(b, input[:]...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)
	return base58.Encode(b)
}

func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

func hash(message []byte) []byte {
	padded := []byte("\x18Bitcoin Signed Message:\n")
	padded = append(padded, numToVarInt(len(message))...)
	padded = append(padded, message...)

	h := sha256.Sum256(padded)
	h2 := sha256.Sum256(h[:])
	return h2[:]
}

func coincurveSig(sign []byte) ([]byte, error) {
	if len(sign) != 65 {
		return nil, fmt.Errorf("invalid sign length: %d", len(sign))
	}
	recoveryID := (sign[0] - 27) & 3
	if recoveryID >= 4 {
		return nil, fmt.Errorf("recovery ID %d not supported", recoveryID)
	}

	return append(sign[1:], recoveryID), nil
}

// https://github.com/vbuterin/pybitcointools/blob/87806f3c984e258a5f30814a089b5c29cbcf0952/bitcoin/main.py#L397
func numToVarInt(n int) []byte {
	const base = 256

	var result []byte
	var minLength int
	if n < 253 {
		return []byte{byte(n)}
	} else if n < 65536 {
		minLength = 2
		result = append(result, 253)
	} else if n < 4294967296 {
		minLength = 4
		result = append(result, 254)
	} else {
		minLength = 8
		result = append(result, 255)
	}

	for n > 0 {
		index := n % base
		result = append(result, byte(index))
		n /= base
	}

	for len(result) <= minLength {
		result = append(result, 0)
	}

	return result
}
