package uiserver

import "testing"

func Test_decodeCmd(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantCmd string
		wantID  int
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
			got, err := decodeCmd(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.CMD != tt.wantCmd {
				t.Errorf("decodeCmd() = %v, want %v", got, tt.wantCmd)
			}
			if got.ID != tt.wantID {
				t.Errorf("decodeCmd() = %v, want %v", got, tt.wantID)
			}
		})
	}
}
