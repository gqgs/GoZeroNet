//go:build !(!go1.15 || go1.17)
// +build go1.15,!go1.17

package serialize

import "github.com/bytedance/sonic"

// Sonic doesn't have 1.17 support yet

var JSONUnmarshal = sonic.Unmarshal
var JSONMarshal = sonic.Marshal
