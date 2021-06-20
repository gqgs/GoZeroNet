package storage

import (
	"database/sql"
	"io"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	io.Closer
	Begin() (Transaction, error)
}

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Rollback() error
	Commit() error
}

type sqliteStorage struct {
	db *sql.DB
}

func (s *sqliteStorage) Begin() (Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func NewStorage(dbName string) (*sqliteStorage, error) {
	db, err := sql.Open("sqlite3", dbName+"?_synchronous=off&cache=shared")
	if err != nil {
		return nil, err
	}
	return &sqliteStorage{
		db: db,
	}, nil
}

func (s *sqliteStorage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
