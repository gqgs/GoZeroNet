package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/ripemd160"
)

type Encoding int

const (
	Hex Encoding = iota
	Base58
)

// PrivateKeyToAddress receives an encoded private key and returns
// the address derived from the associated public key.
// It guesses the key encoding based on the input length.
func PrivateKeyToAddress(encodedKey string) (string, error) {
	var keyBytes []byte
	var err error
	if len(encodedKey) == 64 {
		keyBytes, err = hex.DecodeString(encodedKey)
	} else {
		keyBytes, _, err = base58.CheckDecode(encodedKey)
	}
	if err != nil {
		return "", err
	}
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), keyBytes)
	return PublicKeyToAddress(pubKey.SerializeUncompressed()), nil
}

// NewPrivateKey returns a new random private key encoded in the requested format.
// It panics if it fails to generate the key.
func NewPrivateKey(encoding Encoding) string {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		panic(err)
	}
	switch encoding {
	case Hex:
		return hex.EncodeToString(privKey.Serialize())
	case Base58:
		return base58CheckEncode(privKey.Serialize())
	default:
		panic("invalid encoding")
	}
}

// AuthPrivateKey returns a deterministic base58 encoded key derived in function
// of the seed and the given address.
func AuthPrivateKey(seed string, address string) (string, error) {
	bigInt := new(big.Int)
	bigInt.SetString(hex.EncodeToString([]byte(address)), 16)
	div := new(big.Int)
	div.SetInt64(100_000_000)
	res := bigInt.Mod(bigInt, div)
	child := uint32(res.Uint64())

	master, err := hdkeychain.NewMaster([]byte(seed), &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	derivedChild, err := master.Child(child)
	if err != nil {
		return "", err
	}

	privKey, err := derivedChild.ECPrivKey()
	if err != nil {
		return "", err
	}

	key := append([]byte{0x80}, privKey.Serialize()...)
	return base58CheckEncode(key), nil
}

// IsValidSignature return true if message was signed with key related to address.
func IsValidSignature(message []byte, base64Sign, address string) bool {
	pubkey, err := RecoverPublicKey(message, base64Sign)
	if err != nil {
		return false
	}
	return address == PublicKeyToAddress(pubkey)
}

// RecoverPublicKey recovers the public key given a message and the message signature.
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

// PublicKeyToAddress converts a public key to ZN style addresses.
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
// which is a different behavior from the Python lib.
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

// It returns the ID ZN uses for optional files
func HashID(hexDigest string) (int, error) {
	return hashID(hexDigest, 4)
}

// Given a hex string it returns the first `length` character in base 10.
func hashID(hexDigest string, length int) (int, error) {
	if len(hexDigest) < length {
		return 0, fmt.Errorf("input is too small: %d < %d", len(hexDigest), length)
	}

	res, err := strconv.ParseInt(hexDigest[:length], 16, 0)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}
