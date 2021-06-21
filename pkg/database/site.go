package database

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

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

func (d *siteDatabase) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	return d.storage.Query(query, args...)
}

func (d *siteDatabase) Rebuild() error {
	schema, err := loadDBSchemaFromFile(d.site)
	if err != nil {
		return err
	}

	if err := d.rebuild(schema); err != nil {
		return err
	}

	return d.populate(schema)
}

func (d *siteDatabase) populate(schema *Schema) error {
	regexFunc := make(map[*regexp.Regexp]Map)
	for expr, tableMap := range schema.Maps {
		regex, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("regex error (%s): %s ", err, expr)
		}
		regexFunc[regex] = tableMap
	}

	tx, err := d.storage.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	root := path.Join(config.DataDir, d.site)
	err = filepath.WalkDir(root, func(innerPath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		for regexFunc, tableMap := range regexFunc {
			if regexFunc.MatchString(innerPath) {
				if err := tableMap.ProcessFile(innerPath, tx); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *siteDatabase) rebuild(schema *Schema) error {
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

	// SQLite doesn't like when an open database is moved ("attempt to write a readonly database").
	// So close the database rename it to the correct path and then open the connection again.
	if err := d.storage.Close(); err != nil {
		return err
	}
	if err := newStorage.Close(); err != nil {
		return err
	}
	dbPath := path.Join(config.DataDir, d.site, safe.CleanPath(schema.DBFile))
	if err := os.Rename(tempFile.Name(), dbPath); err != nil {
		return err
	}
	newStorage, err = storage.NewStorage(dbPath)
	d.storage = newStorage
	return err
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
