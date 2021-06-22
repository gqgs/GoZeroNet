package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gqgs/go-zeronet/pkg/lib/storage"
	simdjson "github.com/minio/simdjson-go"
)

var supportsSIMD bool

func init() {
	supportsSIMD = simdjson.SupportedCPU()
}

type Schema struct {
	DBName  string           `json:"db_name"`
	DBFile  string           `json:"db_file"`
	Version int              `json:"version"`
	Maps    map[string]Map   `json:"maps"`
	Tables  map[string]Table `json:"tables"`
}

type Map struct {
	ToTable     []string `json:"to_table"`
	ToJSONTable []string `json:"to_json_table"`
	ToKeyValue  []string `json:"to_keyvalue"`
	FileName    string   `json:"file_name"`
}

type Table struct {
	Cols          [][]string `json:"cols"`
	Indexes       []string   `json:"indexes"`
	SchemaChanged int        `json:"schema_changed"`
}

func getDirectoryAndFilename(innerPath string) (dir string, filename string, err error) {
	i := strings.Index(innerPath, "users/")
	if i < 0 {
		return "", "", errors.New("unexpected path format")
	}

	filename = filepath.Base(innerPath)
	dir = strings.TrimSuffix(innerPath[i:], "/"+filename)
	return
}

func (m *Map) ProcessFile(innerPath string, tx storage.Transaction) error {
	file, err := os.ReadFile(innerPath)
	if err != nil {
		return err
	}

	if !supportsSIMD {
		// TODO: fallback to iter JSON?
		return errors.New("CPU not supported by simdjson")
	}

	dir, filename, err := getDirectoryAndFilename(innerPath)
	if err != nil {
		return err
	}

	if m.FileName != "" {
		filename = m.FileName
	}

	jsonRow, err := tx.Exec(`
		INSERT INTO json (directory, file_name)
		VALUES (?, ?)
		ON CONFLICT(directory, file_name) DO UPDATE SET
		directory = excluded.directory,
		file_name =  excluded.file_name`,
		dir, filename)
	if err != nil {
		return err
	}

	jsonRowID, err := jsonRow.LastInsertId()
	if err != nil {
		return err
	}

	parser, err := simdjson.Parse(file, nil)
	if err != nil {
		return err
	}

	iter := parser.Iter()

	typ := iter.Advance()
	if typ != simdjson.TypeRoot {
		return errors.New("file must start with root element")
	}

	tmp := &simdjson.Iter{}
	if typ, tmp, err = iter.Root(tmp); err != nil {
		return err
	}

	if typ != simdjson.TypeObject {
		return errors.New("root must start an object element")
	}

	obj, err := tmp.Object(nil)
	if err != nil {
		return err
	}

	rootElements, err := obj.Parse(nil)
	if err != nil {
		return err
	}

	tables := make(map[string]struct{})
	for _, table := range m.ToTable {
		tables[table] = struct{}{}
	}

	jsonFields := make(map[string]struct{})
	for _, field := range m.ToJSONTable {
		jsonFields[field] = struct{}{}
	}

	// TODO: refactor this
	for _, rootElement := range rootElements.Elements {
		//nolint:exhaustive
		switch rootElement.Type {
		case simdjson.TypeString:
			if _, isJSONField := jsonFields[rootElement.Name]; isJSONField {
				str, err := rootElement.Iter.String()
				if err != nil {
					return err
				}
				update := fmt.Sprintf("UPDATE json SET %s = ? WHERE rowid = ?", rootElement.Name)
				if _, err := tx.Exec(update, str, jsonRowID); err != nil {
					return err
				}
			}
		case simdjson.TypeFloat:
			if _, isJSONField := jsonFields[rootElement.Name]; isJSONField {
				value, err := rootElement.Iter.Float()
				if err != nil {
					return err
				}
				update := fmt.Sprintf("UPDATE json SET %s = ? WHERE rowid = ?", rootElement.Name)
				if _, err := tx.Exec(update, value, jsonRowID); err != nil {
					return err
				}
			}
		case simdjson.TypeInt:
			if _, isJSONField := jsonFields[rootElement.Name]; isJSONField {
				value, err := rootElement.Iter.Int()
				if err != nil {
					return err
				}
				update := fmt.Sprintf("UPDATE json SET %s = ? WHERE rowid = ?", rootElement.Name)
				if _, err := tx.Exec(update, value, jsonRowID); err != nil {
					return err
				}
			}
		case simdjson.TypeArray:
			newIter := rootElement.Iter
			typ := newIter.Advance()
			if typ == simdjson.TypeObject {
				obj, err = newIter.Object(obj)
				if err != nil {
					return err
				}

				elements, err := obj.Parse(nil)
				if err != nil {
					return err
				}

				var cols []string
				var holders []string
				var values []interface{}
				for _, element := range elements.Elements {
					cols = append(cols, element.Name)
					holders = append(holders, "?")
					switch element.Type {
					case simdjson.TypeString:
						value, err := element.Iter.String()
						if err != nil {
							return err
						}
						values = append(values, value)
					case simdjson.TypeFloat:
						value, err := element.Iter.Float()
						if err != nil {
							return err
						}
						values = append(values, value)
					case simdjson.TypeInt, simdjson.TypeUint:
						value, err := element.Iter.Int()
						if err != nil {
							return err
						}
						values = append(values, value)
					case simdjson.TypeBool:
						value, err := element.Iter.Bool()
						if err != nil {
							return err
						}
						values = append(values, value)
					case simdjson.TypeNull, simdjson.TypeArray, simdjson.TypeObject:
						values = append(values, nil)
					case simdjson.TypeRoot, simdjson.TypeNone:
						return errors.New("invalid type")
					}
				}

				if _, hasTable := tables[rootElement.Name]; hasTable {
					colString := strings.Join(cols, ",")
					holderString := strings.Join(holders, ",")
					query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, rootElement.Name, colString, holderString)
					result, err := tx.Exec(query, values...)
					if err != nil {
						return err
					}
					id, err := result.LastInsertId()
					if err != nil {
						return err
					}
					update := fmt.Sprintf("UPDATE OR IGNORE %s SET json_id = %d WHERE rowid = %d", rootElement.Name, jsonRowID, id)
					if _, err := tx.Exec(update); err != nil {
						return err
					}
				}
			}
		case simdjson.TypeObject, simdjson.TypeNull:
			// can't map these types => no-op
		default:
			return fmt.Errorf("unexpected type: %s", rootElement.Type)
		}
	}

	return nil
}

func (s *Schema) Queries() []string {
	queries := []string{
		"CREATE TABLE _version (key TEXT, value INTEGER)",
		`CREATE UNIQUE INDEX keyindex ON _version(key)`,
		fmt.Sprintf(`INSERT INTO _version VALUES ("db", %d)`, s.Version),
	}

	for tableName, table := range s.Tables {
		cols := make([]string, len(table.Cols))
		for i, col := range table.Cols {
			cols[i] = strings.Join(col, " ")
		}
		queries = append(queries, fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(cols, ",")))
		queries = append(queries, table.Indexes...)
		queries = append(queries, fmt.Sprintf("INSERT INTO _version VALUES (%q, %d)", tableName, table.SchemaChanged))
	}

	return queries
}
