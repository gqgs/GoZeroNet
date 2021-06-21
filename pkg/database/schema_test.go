package database

import (
	"encoding/json"
	"os"
	"testing"
)

func Test_SchemaQueries(t *testing.T) {
	tests := []struct {
		name     string
		dbschema string
		expected []string
	}{
		{
			"given a regular schema it should generate all the queries",
			"testdata/mc-dbschema.json",
			[]string{
				"CREATE TABLE _version (key TEXT, value INTEGER)",
				"CREATE UNIQUE INDEX keyindex ON _version(key)",
				`INSERT INTO _version VALUES ("db", 2)`,
				`INSERT INTO _version VALUES ("posts", 15)`,
				`INSERT INTO _version VALUES ("boards", 12)`,
				`INSERT INTO _version VALUES ("modlogs", 13)`,
				`INSERT INTO _version VALUES ("json", 13)`,

				"CREATE TABLE posts (id TEXT,uri TEXT,thread INTEGER,subject TEXT,body TEXT,username TEXT,time INTEGER,files TEXT,directory TEXT,last_edited INTEGER,capcode INTEGER,json_id INTEGER REFERENCES json (json_id))",
				"CREATE TABLE boards (uri TEXT,title TEXT,config TEXT,description TEXT,json_id INTEGER REFERENCES json (json_id))",
				"CREATE TABLE modlogs (uri TEXT,time INTEGER,action INTEGER,info TEXT,json_id INTEGER REFERENCES json (json_id))",
				"CREATE TABLE json (json_id INTEGER PRIMARY KEY AUTOINCREMENT,directory TEXT,file_name TEXT,cert_user_id TEXT)",

				"CREATE UNIQUE INDEX uuid ON posts(id)",
				"CREATE UNIQUE INDEX post_id ON posts(id,json_id)",
				"CREATE INDEX post_thread ON posts(thread)",
				"CREATE INDEX post_uri ON posts(uri)",
				"CREATE INDEX post_directory ON posts(directory)",
				"CREATE INDEX post_uri_directory ON posts(uri,directory)",
				"CREATE UNIQUE INDEX board_id ON boards(uri,json_id)",
				"CREATE INDEX board_uri ON boards(uri)",
				"CREATE UNIQUE INDEX log_id ON modlogs(uri,action,info,json_id)",
				"CREATE INDEX modlogs_action ON modlogs(action)",
				"CREATE INDEX modlogs_info ON modlogs(info)",
				"CREATE UNIQUE INDEX path ON json(directory, file_name)",
				"CREATE INDEX json_directory ON json(directory)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.dbschema)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			schema := new(Schema)
			if err := json.NewDecoder(file).Decode(schema); err != nil {
				t.Fatal(err)
			}

			expected := make(map[string]struct{}, len(tt.expected))
			for _, query := range tt.expected {
				expected[query] = struct{}{}
			}

			created := schema.Queries()

			for _, c := range created {
				t.Log(c)
				delete(expected, c)
			}

			for e := range expected {
				t.Errorf("query not created: %q", e)
			}
		})
	}
}

func Test_getDirectoryAndFilename(t *testing.T) {
	tests := []struct {
		name          string
		innerPath     string
		wantDirectory string
		wantPath      string
	}{
		{
			"it should extract dir and data.json",
			"data/1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D/data/users/1LjoRHqpAi5FFokP3aYUJuavw7c3XS7HVY/data.json",
			"users/1LjoRHqpAi5FFokP3aYUJuavw7c3XS7HVY",
			"data.json",
		},
		{
			"it should extract dir and content.json",
			"data/1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D/data/users/1LjoRHqpAi5FFokP3aYUJuavw7c3XS7HVY/content.json",
			"users/1LjoRHqpAi5FFokP3aYUJuavw7c3XS7HVY",
			"content.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _ := getDirectoryAndFilename(tt.innerPath)
			if got != tt.wantDirectory {
				t.Errorf("getDirectoryAndFilename() got = %v, want dir %v", got, tt.wantDirectory)
			}
			if got1 != tt.wantPath {
				t.Errorf("getDirectoryAndFilename() got = %v, want path %v", got1, tt.wantPath)
			}
		})
	}
}
