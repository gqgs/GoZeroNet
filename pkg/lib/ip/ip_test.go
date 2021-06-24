package ip

import (
	"encoding/binary"
	"testing"
)

func TestParseIPv4(t *testing.T) {
	tests := []struct {
		name string
		addr []byte
		want string
	}{
		{
			"valid ipv4",
			[]byte{0xAD, 0x2F, 0x7B, 0x6F, 0x04, 0x3C},
			"173.47.123.111:15364",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseIPv4(tt.addr, binary.LittleEndian); got != tt.want {
				t.Errorf("ParseIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}
