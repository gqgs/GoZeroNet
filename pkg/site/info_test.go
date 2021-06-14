package site

import (
	"testing"
)

func Test_addressShort(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{
			"valid address",
			"1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D",
			"1HeLLo..Tf3D",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addressShort(tt.addr); got != tt.want {
				t.Errorf("addressShort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressHash(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{
			"valid address",
			"1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D",
			"f69941233e191d9e00f0cd16c5da10b0124d1c0a498b5ecfa1448b21a3eb0094",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addressHash(tt.addr); got != tt.want {
				t.Errorf("addressHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
