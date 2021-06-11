package uiserver

import "testing"

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
