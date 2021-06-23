package safe

import (
	"os"
	"path/filepath"
	"sync/atomic"
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

// Counter returns a closure that will generate sequential IDs.
// It it safe to call the enclosed function from multiple goroutines concurrently.
func Counter() func() int64 {
	var id int64
	return func() int64 {
		return atomic.AddInt64(&id, 1)
	}
}
