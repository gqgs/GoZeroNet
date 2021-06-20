package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/lib/storage"
)

type siteDatabase struct {
	site    string
	storage storage.Storage
}

func NewSiteDatabase(site string) (*siteDatabase, error) {
	schema, err := loadDBSchemaFromFile(site)
	if err != nil {
		return nil, err
	}

	dbPath := path.Join(config.DataDir, site, safe.CleanPath(schema.DBFile))
	storage, err := storage.NewStorage(dbPath)
	if err != nil {
		return nil, err
	}

	return &siteDatabase{
		site:    site,
		storage: storage,
	}, nil
}

func (d *siteDatabase) Close() error {
	if d == nil {
		return nil
	}
	return d.storage.Close()
}

func (d *siteDatabase) Rebuild() error {
	schema, err := loadDBSchemaFromFile(d.site)
	if err != nil {
		return err
	}

	dbDir := path.Dir(path.Join(config.DataDir, d.site, safe.CleanPath(schema.DBFile)))
	tempFile, err := ioutil.TempFile(dbDir, "")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	newStorage, err := storage.NewStorage(tempFile.Name())
	if err != nil {
		return err
	}

	tx, err := newStorage.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, query := range schema.Queries() {
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("%s: %q", err, query)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	// replace old storage with new one
	d.storage.Close()
	d.storage = newStorage

	dbPath := path.Join(config.DataDir, d.site, safe.CleanPath(schema.DBFile))

	return os.Rename(tempFile.Name(), dbPath)
}

func loadDBSchemaFromFile(site string) (*Schema, error) {
	file, err := os.Open(path.Join(config.DataDir, site, "dbschema.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	schema := new(Schema)
	return schema, json.NewDecoder(file).Decode(schema)
}
