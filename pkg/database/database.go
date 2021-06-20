package database

import (
	"io"
)

type Database interface {
	io.Closer
	Rebuild() error
}
