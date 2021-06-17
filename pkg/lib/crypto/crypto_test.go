package crypto

import (
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
			addr, err := PublicKeyToAddress(pubKey)
			require.NoError(t, err)
			require.Equal(t, tt.wantAddress, addr)
		})
	}
}
