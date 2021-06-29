package site

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sortDownloads(t *testing.T) {
	files := map[string]File{
		"README.md":                   {},
		"LICENSE":                     {},
		"chartjs/chart.bundle.min.js": {},
		"index.html":                  {},
		"languages/ru.json":           {},
		"dbschema.json":               {},
	}
	sorted := sortDownloads(files)
	require.Len(t, sorted, len(files))
	require.Equal(t, "dbschema.json", sorted[0])
	require.Equal(t, "index.html", sorted[1])
}
