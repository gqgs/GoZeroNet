package site

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Download(t *testing.T) {
	contentFile, err := os.Open("./testdata/content.json")
	if err != nil {
		t.Fatal(err)
	}
	defer contentFile.Close()

	contentData, err := io.ReadAll(contentFile)
	if err != nil {
		t.Fatal(err)
	}

	content := new(Content)
	if err := json.Unmarshal(contentData, content); err != nil {
		t.Fatal(err)
	}

	require.True(t, content.isValid())
}
