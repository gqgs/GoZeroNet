package safe

import (
	"os"
	"path/filepath"
)

// CleanPath makes a path safe for use with filepath.Join.
func CleanPath(path string) string {
	if path == "" {
		return ""
	}

	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		path = filepath.Clean(string(os.PathSeparator) + path)
		path, _ = filepath.Rel(string(os.PathSeparator), path)
	}
	return filepath.Clean(path)
}
