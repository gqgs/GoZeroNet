//go:build !go1.15 || go1.17
// +build !go1.15 go1.17

package serialize

import "encoding/json"

var JSONUnmarshal = json.Unmarshal
var JSONMarshal = json.Marshal
