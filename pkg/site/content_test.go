package site

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isValid(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		wantIsValid bool
	}{
		{
			"root content.json",
			"./testdata/root-content.json",
			true,
		},
		{
			"data user content.json",
			"./testdata/datauser-content.json",
			true,
		},
		{
			"user content.json",
			"./testdata/user-content.json",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFile, err := os.Open(tt.filename)
			if err != nil {
				t.Error(err)
				return
			}
			defer contentFile.Close()

			contentData, err := io.ReadAll(contentFile)
			if err != nil {
				t.Error(err)
				return
			}

			content := new(Content)
			if err := json.Unmarshal(contentData, content); err != nil {
				t.Error(err)
				return
			}

			require.Equal(t, tt.wantIsValid, content.isValid())
		})
	}
}
