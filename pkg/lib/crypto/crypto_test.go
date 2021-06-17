package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecoverPublicKey(t *testing.T) {
	tests := []struct {
		name        string
		message     []byte
		signature   string
		wantAddress string
	}{
		{
			"given a valid signature it should return the public key address",
			[]byte("1:1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D"),
			"HLcq242ZHh4nTexhe6kvkBroycZ1JpF4pjlLGxbhjKAwDAfdCZ/gxUwM9aIN6OrD8K5YqAfvIVlbwkLMB1XSEDo=",
			"1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubKey, err := RecoverPublicKey(tt.message, tt.signature)
			require.NoError(t, err)
			addr := PublicKeyToAddress(pubKey)
			require.Equal(t, tt.wantAddress, addr)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numToVarInt(tt.n); tt.hexWant != hex.EncodeToString(got) {
				t.Errorf("numToVarInt() = %x, want %s", got, tt.hexWant)
			}
		})
	}
}
