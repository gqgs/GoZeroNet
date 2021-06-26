package storage

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"sync"

	"github.com/gqgs/go-zeronet/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

type execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Storage interface {
	io.Closer
	execer
	queryer
	Begin() (Transaction, error)
	// QueryObjectList is a helper to construct and return an object list of generic type.
	QueryObjectList(query string, args ...interface{}) ([]map[string]interface{}, error)
}

type Transaction interface {
	driver.Tx
	execer
	queryer
}

type sqliteStorage struct {
	mu sync.RWMutex
	db *sql.DB
}

func (s *sqliteStorage) Query(query string, args ...interface{}) (*sql.Rows, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Query(query, args...)
}

func (s *sqliteStorage) Exec(query string, args ...interface{}) (sql.Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Exec(query, args...)
}

func (s *sqliteStorage) QueryObjectList(query string, args ...interface{}) ([]map[string]interface{}, error) {
	s.mu.RLock()
	rows, err := s.db.Query(query, args...)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if config.ValidateDatabaseQueries {
		columnsMap := make(map[string]struct{}, len(columns))
		for _, c := range columns {
			if _, found := columnsMap[c]; found {
				return nil, fmt.Errorf("ambiguous column name: %q", c)
			}
			columnsMap[c] = struct{}{}
		}
	}

	values := make([]interface{}, len(columns))
	rowValues := make([]interface{}, len(columns))
	for i := range values {
		rowValues[i] = &values[i]
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		if err := rows.Scan(rowValues...); err != nil {
			return nil, err
		}

		rowResult := make(map[string]interface{})

		for i, value := range values {
			switch value.(type) {
			case nil:
			default:
				rowResult[columns[i]] = value
			}
		}

		if len(rowResult) == 0 {
			continue
		}

		result = append(result, rowResult)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return result, nil
}

func (s *sqliteStorage) Begin() (Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Begin()
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
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Close()
}
