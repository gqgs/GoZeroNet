package database

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/storage"
)

var ErrFileNotFound = fmt.Errorf("content: %w", os.ErrNotExist)

type Info interface {
	Update(site string, broadcaster event.Broadcaster)
	GetSize() int
	GetIsDownloaded() bool
	AddUploaded(uploaded int)
}

type ContentDatabase interface {
	io.Closer
	UpdateFile(site string, fileInfo *event.FileInfo) error
	UpdatePeer(site string, peerInfo *event.PeerInfo) error
	UpdateContent(site string, contentInfo *event.ContentInfo) error
	FileInfo(site, innerPath string) (*event.FileInfo, error)
	UpdatedFiles(site string, since time.Time) ([]string, error)
	UpdatedContent(site string, since time.Time) (map[string]int, error)
	Peers(site string, need int) ([]string, error)
	ContentInfo(site, innerPath string) (*event.ContentInfo, error)
	Info(site, innerPath string) (Info, error)
}

type contentDatabase struct {
	storage storage.Storage
}

func (c *contentDatabase) Info(site, innerPath string) (Info, error) {
	if strings.HasSuffix(innerPath, "content.json") {
		return c.ContentInfo(site, innerPath)
	}
	return c.FileInfo(site, innerPath)
}

func (c *contentDatabase) UpdatedFiles(site string, since time.Time) ([]string, error) {
	const query = `
		SELECT f.inner_path FROM file f INNER JOIN site s USING(site_id)
		WHERE s.address = ? AND f.time_added >= ? AND f.downloaded = f.size
		UNION
		SELECT c.inner_path FROM content c INNER JOIN SITE s USING(site_id)
		WHERE s.address = ? AND c.time_added >= ?
	`
	rows, err := c.storage.Query(query, site, since.UTC(), site, since.UTC())
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

	return files, rows.Err()
}

func (c *contentDatabase) FileInfo(site, innerPath string) (*event.FileInfo, error) {
	query := `
		SELECT f.inner_path, f.hash, f.size, f.downloaded = f.size, f.is_pinned, f.is_optional, f.uploaded, f.piece_size, f.piecemap, f.downloaded, (f.downloaded * 1.0 / f.size) * 100, f.peer
		FROM file f INNER JOIN site s USING(site_id)
		WHERE f.inner_path = ? AND s.address = ?
	`
	rows, err := c.storage.Query(query, innerPath, site)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := new(event.FileInfo)
	if rows.Next() {
		if err := rows.Scan(
			&info.InnerPath,
			&info.Hash,
			&info.Size,
			&info.IsDownloaded,
			&info.IsPinned,
			&info.IsOptional,
			&info.Uploaded,
			&info.PieceSize,
			&info.Piecemap,
			&info.Downloaded,
			&info.DownloadedPercent,
			&info.Peer,
		); err != nil {
			return nil, err
		}
		return info, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return info, fmt.Errorf("%w: %s", ErrFileNotFound, innerPath)
}

func (c *contentDatabase) Close() error {
	if c == nil || c.storage == nil {
		return nil
	}
	return c.storage.Close()
}

func (c *contentDatabase) UpdateFile(site string, info *event.FileInfo) error {
	tx, err := c.storage.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("INSERT OR IGNORE INTO site (address) VALUES (?)", site); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO file (site_id, inner_path, hash, size, is_pinned, is_optional, uploaded, piece_size, piecemap, downloaded, peer)
		VALUES ((SELECT site_id FROM site WHERE address = ?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (site_id, inner_path) DO
		UPDATE SET
			is_pinned = excluded.is_pinned,
			is_optional = excluded.is_optional,
			uploaded = excluded.uploaded,
			piece_size = excluded.piece_size,
			piecemap = excluded.piecemap,
			time_added = CURRENT_TIMESTAMP,
			hash = excluded.hash,
			size = excluded.size,
			downloaded = excluded.downloaded,
			peer = excluded.peer
		`,
		site,
		info.InnerPath,
		info.Hash,
		info.Size,
		info.IsPinned,
		info.IsOptional,
		info.Uploaded,
		info.PieceSize,
		info.Piecemap,
		info.Downloaded,
		info.Peer,
	); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *contentDatabase) UpdatePeer(site string, peerInfo *event.PeerInfo) error {
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
		site, peerInfo.Address, peerInfo.ReputationDelta); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *contentDatabase) UpdateContent(site string, info *event.ContentInfo) error {
	tx, err := c.storage.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("INSERT OR IGNORE INTO site (address) VALUES (?)", site); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO content (site_id, inner_path, modified, size)
		VALUES ((SELECT site_id FROM site WHERE address = ?), ?, ?, ?)
		ON CONFLICT (site_id, inner_path) DO
		UPDATE SET
			modified = excluded.modified,
			size = excluded.size,
			time_added = CURRENT_TIMESTAMP
		`,
		site,
		info.InnerPath,
		info.Modified,
		info.Size,
	); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *contentDatabase) UpdatedContent(site string, since time.Time) (map[string]int, error) {
	query := `
		SELECT c.inner_path, c.modified
		FROM content c INNER JOIN site s USING(site_id)
		WHERE s.address = ? AND c.modified > ?
		ORDER BY c.modified DESC
		LIMIT 100
	`
	rows, err := c.storage.Query(query, site, since.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contentList := make(map[string]int)
	if rows.Next() {
		var innerPath string
		var modified int
		if err := rows.Scan(&innerPath, &modified); err != nil {
			return nil, err
		}
		contentList[innerPath] = modified
	}
	return contentList, rows.Err()

}

func (c *contentDatabase) Peers(site string, need int) ([]string, error) {
	var peers []string
	query := `
		SELECT p.address
		FROM peer p INNER JOIN site s USING(site_id)
		WHERE s.address = ? AND p.reputation > 0
		ORDER BY p.time_added DESC
		LIMIT ?
	`
	rows, err := c.storage.Query(query, site, need)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peer string
	for rows.Next() {
		if err := rows.Scan(&peer); err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}
	return peers, rows.Err()
}

func (c *contentDatabase) ContentInfo(site, innerPath string) (*event.ContentInfo, error) {
	query := `
		SELECT c.inner_path, c.modified, c.size
		FROM content c INNER JOIN site s USING(site_id)
		WHERE s.address = ? and c.inner_path = ?
	`
	rows, err := c.storage.Query(query, site, innerPath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := new(event.ContentInfo)
	if rows.Next() {
		if err := rows.Scan(&info.InnerPath, &info.Modified, &info.Size); err != nil {
			return nil, err
		}
		return info, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return info, fmt.Errorf("%w: %s", ErrFileNotFound, innerPath)
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
		`CREATE TABLE IF NOT EXISTS file (file_id INTEGER PRIMARY KEY UNIQUE NOT NULL, site_id INTEGER NOT NULL REFERENCES site (site_id) ON DELETE CASCADE, inner_path TEXT, hash TEXT, size INTEGER, peer INTEGER DEFAULT 0, uploaded INTEGER DEFAULT 0, downloaded INTEGER DEFAULT 0, is_pinned INTEGER DEFAULT 0, is_optional INTEGER DEFAULT 0, time_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP, piece_size INTEGER, piecemap TEXT)`,
		`CREATE INDEX IF NOT EXISTS file_path ON file (inner_path)`,
		`CREATE INDEX IF NOT EXISTS file_hash ON file (hash)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS file_path_hash ON file (site_id, inner_path)`,

		// Content
		`CREATE TABLE IF NOT EXISTS content (content_id INTEGER PRIMARY KEY UNIQUE NOT NULL, site_id INTEGER REFERENCES site (site_id) ON DELETE CASCADE, inner_path TEXT, modified INTEGER, size INTEGER, time_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP, time_downloaded INTEGER DEFAULT 0, time_accessed INTEGER DEFAULT 0)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS content_key ON content (site_id, inner_path)`,
		`CREATE INDEX IF NOT EXISTS content_modified ON content (site_id, modified)`,
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
