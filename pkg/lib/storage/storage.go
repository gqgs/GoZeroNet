package storage

import (
	"database/sql"
	"database/sql/driver"
	"io"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	io.Closer
	Begin() (Transaction, error)
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)
}

type Transaction interface {
	driver.Tx
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type sqliteStorage struct {
	db *sql.DB
}

func (s *sqliteStorage) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
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
