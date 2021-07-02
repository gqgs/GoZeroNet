package site

import (
	"time"

	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
)

// OpenDB opens a new connection to the site's database.
// The caller is responsible for calling Close when
// the database is no longer needed
func (s *Site) OpenDB() error {
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
	return s.db.Close()
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

func (s *Site) UpdateDB(since time.Time) error {
	updated, err := s.contentDB.UpdatedFiles(s.addr, since)
	if err != nil {
		return err
	}
	s.log.WithField("updated", len(updated)).Info("updating database")
	if len(updated) == 0 {
		return nil
	}
	return s.db.Update(updated...)
}
