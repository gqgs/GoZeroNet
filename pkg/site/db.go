package site

import (
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
)

// OpenDB opens a new connection to the site's database.
// The caller is responsible for calling Close when
// the database is no longer needed
func (s *Site) OpenDB() error {
	if s.db != nil {
		return nil
	}

	db, err := database.NewSiteDatabase(s.addr)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Site) CloseDB() error {
	if s == nil || s.db == nil {
		return nil
	}
	db := s.db
	s.db = nil
	return db.Close()
}

func (s *Site) RebuildDB() error {
	return s.db.Rebuild()
}

func (s *Site) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	return s.db.Query(query, args...)
}

func (s *Site) FileInfo(innerPath string) (*event.FileInfo, error) {
	return s.contentDB.FileInfo(s.addr, innerPath)
}

func (s *Site) hasDB() bool {
	if s.db == nil {
		return false
	}
	return s.db.Exists()
}
