package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidSignature(t *testing.T) {
	tests := []struct {
		name      string
		message   []byte
		signature string
		address   string
		want      bool
	}{
		{
			"given a valid signature it should return the public key address",
			[]byte("1:1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D"),
			"HLcq242ZHh4nTexhe6kvkBroycZ1JpF4pjlLGxbhjKAwDAfdCZ/gxUwM9aIN6OrD8K5YqAfvIVlbwkLMB1XSEDo=",
			"1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsValidSignature(tt.message, tt.signature, tt.address))
		})
	}
}

func Test_numToVarInt(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		hexWant string
	}{
		{
			"n < 253",
			200,
			"c8",
		},
		{
			"253 <= n < 65536",
			45678,
			"fd6eb2",
		},
		{
			"65536 <= n < 4294967296",
			4294967290,
			"fefaffffff",
		},
		{
			"4294967296 <= n",
			5194967296,
			"ff00e9a43501000000",
		},
		{
			"n = MAX_INT64",
			9223372036854775807,
			"ffffffffffffffff7f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numToVarInt(tt.n); tt.hexWant != hex.EncodeToString(got) {
				t.Errorf("numToVarInt() = %x, want %s", got, tt.hexWant)
			}
		})
	}
}

func TestPrivateKeyToAddress(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			"given a valid hex encoded key it should return its address",
			"366e9056541340ae10ef5af621d73872b6b678161aa9c0dc409701ca155a6693",
			"1Ea7ZmUiuwNBdYG6v54yyJJfK1NJy9agGX",
		},
		{
			"given a valid base58 encoded key it should return its address",
			"5KRDwnpby7hk3fn2Giov61BTPggwyYqnJSCgopRdprtqqNbgPXo",
			"1HzAPSQjyEDtbQeiWaLvmdu6hSxWcpTwjD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyToAddress(tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewPrivateKey(t *testing.T) {
	tests := []struct {
		name     string
		encoding Encoding
		wantLen  int
	}{
		{
			"hex string",
			Hex,
			64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPrivateKey(tt.encoding)
			require.Len(t, got, tt.wantLen)
		})
	}
}

func Test_AuthPrivateKeyKey(t *testing.T) {
	tests := []struct {
		name        string
		seed        string
		address     string
		wantPrivKey string
	}{
		{
			"given a valid seed it should return a deterministic auth key",
			"e180efa477c63b0f2757eac7b1cce781877177fe0966be62754ffd4c8592ce38",
			"",
			"5JSbeF5PevdrsYjunqpg7kAGbnCVYa1T4APSL3QRu8EoAmXRc7Y",
		},
		{
			"ZeroName seed",
			"366e9056541340ae10ef5af621d73872b6b678161aa9c0dc409701ca155a6693",
			"1Name2NXVi1RDPDgf5617UoW7xA6YrhM9F",
			"5KRDwnpby7hk3fn2Giov61BTPggwyYqnJSCgopRdprtqqNbgPXo",
		},
		{
			"ZeroHello seed",
			"366e9056541340ae10ef5af621d73872b6b678161aa9c0dc409701ca155a6693",
			"1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D",
			"5HwmpV4rm7FFNYzrk1Xy4Gj5WBXeGtbzpZxsBzt11wQsucPnRk7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authKey, err := AuthPrivateKey(tt.seed, tt.address)
			require.NoError(t, err)
			require.Equal(t, tt.wantPrivKey, authKey)
		})
	}
}
func TestHashID(t *testing.T) {
	tests := []struct {
		name      string
		hexDigest string
		want      int
	}{
		{
			"Given a valid hex string it should return the expected hash id",
			"ea2c2acb30bd5e1249021976536574dd3f0fd83340e023bb4e78d0d818adf30a",
			59948,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashID(tt.hexDigest)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
func TestSign(t *testing.T) {
	tests := []struct {
		name          string
		message       []byte
		privateKey    string
		wantSignature string
	}{
		{
			"given a message it should return the expected signature",
			[]byte("1:1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D"),
			"5KRDwnpby7hk3fn2Giov61BTPggwyYqnJSCgopRdprtqqNbgPXo",
			"G6YZwyYQSoPC4VSlFrinM9WXOrCGSsUB6CkXl8Vub0kKUNQmLv6ItuNM4ECeMW9kkJgE1FheqjDM9Ow4cKZ/ZZY=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sign, err := Sign(tt.message, tt.privateKey)
			require.NoError(t, err)
			require.Equal(t, tt.wantSignature, sign)
		})
	}
}
