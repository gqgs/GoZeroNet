package uiwebsocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Benchmark_decode(b *testing.B) {
	const data = `{
			"cmd": "fileGet",
			"params": {
				"inner_path": "data/users/1DJEfWgdJ3rGnqCrKwLrBdQGodRaEBkAry/6e2ae2cc67be7a03005a672aef84f1a9e3cf403a.png",
				"required": true,
				"format": "base64"
			},
			"wrapper_nonce": "719f6898d280b45d89c864021a5f2c8f74f253c08d2f14caba8988de71c11f14",
			"id": 100
		}`
	for i := 0; i < b.N; i++ {
		result, err := decode([]byte(data))
		require.NoError(b, err)
		require.Equal(b, int64(100), result.ID)
		require.Equal(b, "719f6898d280b45d89c864021a5f2c8f74f253c08d2f14caba8988de71c11f14", result.WrapperNonce)
		require.Equal(b, "fileGet", result.CMD)
	}
}
func Test_decode(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantCmd string
		wantID  int64
		wantErr bool
	}{
		{
			"given a json payload it should find the cmd",
			[]byte(`{"cmd":"channelJoin","params":{"channels":["siteChanged","serverChanged"]},"id":1000000}`),
			"channelJoin",
			1000000,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decode(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.CMD != tt.wantCmd {
				t.Errorf("decode() = %v, want %v", got, tt.wantCmd)
			}
			if got.ID != tt.wantID {
				t.Errorf("decode() = %v, want %v", got, tt.wantID)
			}
		})
	}
}
