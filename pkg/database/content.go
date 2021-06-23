package database

import (
	"errors"
	"io"
	"path"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/storage"
)

var ErrFileNotFound = errors.New("file not found")

type ContentDatabase interface {
	io.Closer
	UpdateFile(site, innerPath, hash string, size int) error
	UpdatePeer(site, address string, reputationDelta int) error
	FileInfo(site, innerPath string) (*event.FileInfo, error)
	GetUpdatedFiles(site string, since time.Time) ([]string, error)
}

type contentDatabase struct {
	storage storage.Storage
}

func (c *contentDatabase) GetUpdatedFiles(site string, since time.Time) ([]string, error) {
	const query = "SELECT f.inner_path FROM file f INNER JOIN site s USING(site_id) WHERE s.address = ? AND f.time_added >= ?"
	rows, err := c.storage.Query(query, site, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []string
	var file string
	for rows.Next() {
		if err = rows.Scan(&file); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (c *contentDatabase) FileInfo(site, innerPath string) (*event.FileInfo, error) {
	query := `
		SELECT f.inner_path, f.hash, f.size, f.is_downloaded
		FROM file f INNER JOIN site s USING(site_id)
		WHERE f.inner_path = ? AND s.address = ?
	`
	rows, err := c.storage.Query(query, innerPath, site)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := new(event.FileInfo)
	if next := rows.Next(); next {
		if err := rows.Scan(&info.InnerPath, &info.Hash, &info.Size, &info.IsDownloaded); err != nil {
			return nil, err
		}
		return info, rows.Err()
	}
	return info, ErrFileNotFound
}

func (c *contentDatabase) Close() error {
	if c == nil || c.storage == nil {
		return nil
	}
	return c.storage.Close()
}

func (c *contentDatabase) UpdateFile(site, innerPath, hash string, size int) error {
	tx, err := c.storage.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("INSERT OR IGNORE INTO site (address) VALUES (?)", site); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO file (site_id, inner_path, hash, size)
		VALUES ((SELECT site_id FROM site WHERE address = ?), ?, ?, ?)`,
		site, innerPath, hash, size); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *contentDatabase) UpdatePeer(site, address string, reputationDelta int) error {
	tx, err := c.storage.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("INSERT OR IGNORE INTO site (address) VALUES (?)", site); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO peer (site_id, address, reputation)
		VALUES ((SELECT site_id FROM site WHERE address = ?), ?, ?)
		ON CONFLICT (site_id, address) DO UPDATE SET reputation = reputation + excluded.reputation
		`,
		site, address, reputationDelta); err != nil {
		return err
	}
	return tx.Commit()
}

func NewContentDatabase() (*contentDatabase, error) {
	dbPath := path.Join(config.DataDir, "content.db")
	storage, err := storage.NewStorage(dbPath)
	if err != nil {
		return nil, err
	}

	queries := []string{
		// Site
		`CREATE TABLE IF NOT EXISTS site (site_id INTEGER PRIMARY KEY ASC NOT NULL UNIQUE, address TEXT NOT NULL, time_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS site_address ON site (address)`,

		// Peer
		`CREATE TABLE IF NOT EXISTS peer (site_id INTEGER REFERENCES site (site_id) ON DELETE CASCADE, address TEXT NOT NULL, reputation INTEGER DEFAULT 0, time_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS peer_key ON peer (site_id, address)`,

		// File
		`CREATE TABLE IF NOT EXISTS file (file_id INTEGER PRIMARY KEY UNIQUE NOT NULL, site_id INTEGER NOT NULL REFERENCES site (site_id) ON DELETE CASCADE, inner_path TEXT, hash TEXT, size INTEGER, peer INTEGER DEFAULT 0, uploaded INTEGER DEFAULT 0, is_downloaded INTEGER DEFAULT 1, is_pinned INTEGER DEFAULT 0, time_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`,
		`CREATE INDEX IF NOT EXISTS file_path ON file (inner_path)`,
		`CREATE INDEX IF NOT EXISTS file_hash ON file (hash)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS file_hash ON file (site_id, inner_path, hash)`,
	}

	tx, err := storage.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			return nil, err
		}
	}

	return &contentDatabase{
		storage: storage,
	}, tx.Commit()
}
