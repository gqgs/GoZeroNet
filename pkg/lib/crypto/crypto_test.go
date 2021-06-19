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
		name   string
		hexKey string
		want   string
	}{
		{
			"given a valid hex encoded key it should return its address",
			"366e9056541340ae10ef5af621d73872b6b678161aa9c0dc409701ca155a6693",
			"1Ea7ZmUiuwNBdYG6v54yyJJfK1NJy9agGX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyToAddress(tt.hexKey)
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
