package commands

import (
	"rdbslite/data"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		cmd     Command
		wantErr bool
	}{
		{
			name: "creates table",
			cmd: Command{
				Name: "CREATE",
				Args: []string{"TABLE", "users", "id", "int", "name", "text"},
			},
			wantErr: false,
		},
		{
			name: "rejects empty schema",
			cmd: Command{
				Name: "CREATE",
				Args: []string{"TABLE", "users"},
			},
			wantErr: true,
		},
		{
			name: "rejects invalid schema",
			cmd: Command{
				Name: "CREATE",
				Args: []string{"Table", "users", "int"},
			},
			wantErr: true,
		},
		{
			name: "rejects empty table name",
			cmd: Command{
				Name: "CREATE",
				Args: []string{"Table", " ", "id", "int"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := data.NewDatabase()
			_, err := Create(&db, tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreate_DuplicateTable(t *testing.T) {
	db := data.NewDatabase()
	cmd := Command{
		Name: "CREATE",
		Args: []string{"TABLE", "users", "id", "int"},
	}

	_, err1 := Create(&db, cmd)
	_, err2 := Create(&db, cmd)

	if err1 != nil {
		t.Fatalf("first create failed: %v", err1)
	}
	if err2 == nil {
		t.Fatalf("expected duplicate table error")
	}
}

func TestCreate_DuplicateTableCaseInsensitive(t *testing.T) {
	db := data.NewDatabase()
	cmd := Command{
		Name: "CREATE",
		Args: []string{"TABLE", "users", "id", "int"},
	}

	cmd2 := Command{
		Name: "CREATE",
		Args: []string{"TABLE", "Users", "id", "int"},
	}

	_, err1 := Create(&db, cmd)
	_, err2 := Create(&db, cmd2)

	if err1 != nil {
		t.Fatalf("first create failed: %v", err1)
	}
	if err2 == nil {
		t.Fatalf("expected duplicate table error with case insensitive naming")
	}
}
