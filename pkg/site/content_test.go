package site

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isValid(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			"ZeroHello root content.json",
			"./testdata/0hello-root-content.json",
		},
		{
			"ZeroBlog root content.json",
			"./testdata/0blog-root-content.json",
		},
		{
			"ZeroBlog datauser content.json",
			"./testdata/0blog-datauser-content.json",
		},
		{
			"ZeroBlog user content.json",
			"./testdata/0blog-user-content.json",
		},
		{
			"MC root content.json",
			"./testdata/mc-root-content.json",
		},
		{
			"MC datauser content.json",
			"./testdata/mc-datauser-content.json",
		},
		{
			"MC user content.json",
			"./testdata/mc-user-content.json",
		},
		{
			"0ch root content.json",
			"./testdata/0ch-root-content.json",
		},
		{
			"0ch datauser content.json",
			"./testdata/0ch-datauser-content.json",
		},
		{
			"0ch user content.json",
			"./testdata/0ch-user-content.json",
		},
		{
			"0ch archive content.json",
			"./testdata/0ch-archive-content.json",
		},
		{
			"ZeroTalk root content.json",
			"./testdata/0talk-root-content.json",
		},
		{
			"ZeroTalk datauser content.json",
			"./testdata/0talk-datauser-content.json",
		},
		{
			"ZeroTalk user content.json",
			"./testdata/0talk-user-content.json",
		},
		{
			"MC user content.json with extra file fields",
			"./testdata/mc-user-content-extra-file-fields.json",
		},
		{
			"ZeroUp root content.json",
			"./testdata/0up-root-content.json",
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

			content := new(Content)
			if err := json.NewDecoder(contentFile).Decode(content); err != nil {
				t.Error(err)
				return
			}

			require.True(t, content.isValid())
		})
	}
}
