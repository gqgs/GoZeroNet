package database

import (
	"fmt"
	"strings"
)

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

func (s Schema) Queries() []string {
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
