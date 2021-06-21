package database

import (
	"io"
)

type Database interface {
	io.Closer
	Rebuild() error
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)
}
