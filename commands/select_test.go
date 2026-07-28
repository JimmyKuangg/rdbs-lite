package commands

import (
	"strings"
	"testing"

	"rdbslite/data"
)

func seedSelectDB(t *testing.T) data.Database {
	t.Helper()

	db := data.NewDatabase()

	_, err := Create(&db, Command{
		Name: "CREATE",
		Args: []string{"TABLE", "users", "id", "int", "name", "text", "verified", "bool"},
	})
	if err != nil {
		t.Fatalf("[%s] create failed: %v", t.Name(), err)
	}

	_, err = Insert(&db, Command{
		Name: "INSERT",
		Args: []string{"INTO", "users", "id", "1", "name", "alice", "verified", "true"},
	})
	if err != nil {
		t.Fatalf("[%s] first insert failed: %v", t.Name(), err)
	}

	_, err = Insert(&db, Command{
		Name: "INSERT",
		Args: []string{"INTO", "users", "id", "2", "name", "bob", "verified", "false"},
	})
	if err != nil {
		t.Fatalf("[%s] second insert failed: %v", t.Name(), err)
	}

	return db
}

func TestSelect(t *testing.T) {
	tests := []struct {
		name    string
		cmd     Command
		wantErr bool
	}{
		{
			name: "select wildcard works",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"*", "FROM", "users"},
			},
			wantErr: false,
		},
		{
			name: "select wildcard works with case-insensitive from and table",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"*", "from", "USERS"},
			},
			wantErr: false,
		},
		{
			name: "select named columns works",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"id", "name", "FROM", "users"},
			},
			wantErr: false,
		},
		{
			name: "rejects not enough args",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"*", "FROM"},
			},
			wantErr: true,
		},
		{
			name: "rejects missing FROM keyword",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"id", "name", "users"},
			},
			wantErr: true,
		},
		{
			name: "rejects no columns before FROM",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"FROM", "users", "x"},
			},
			wantErr: true,
		},
		{
			name: "rejects missing table after FROM",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"id", "name", "FROM"},
			},
			wantErr: true,
		},
		{
			name: "rejects unexpected trailing tokens",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"id", "FROM", "users", "WHERE"},
			},
			wantErr: true,
		},
		{
			name: "rejects unknown table",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"*", "FROM", "ghosts"},
			},
			wantErr: true,
		},
		{
			name: "rejects unknown column",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"unknown_col", "FROM", "users"},
			},
			wantErr: true,
		},
		{
			name: "accepts case-insensitive column names",
			cmd: Command{
				Name: "SELECT",
				Args: []string{"ID", "NaMe", "FROM", "users"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := seedSelectDB(t)
			_, err := Select(&db, tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("[%s] err = %v, wantErr = %v", t.Name(), err, tt.wantErr)
			}
		})
	}
}

func TestSelect_ProjectionOutput(t *testing.T) {
	db := seedSelectDB(t)

	out, err := Select(&db, Command{
		Name: "SELECT",
		Args: []string{"name", "id", "FROM", "users"},
	})
	if err != nil {
		t.Fatalf("[%s] select failed: %v", t.Name(), err)
	}

	// Header should include selected columns.
	if !strings.Contains(out, "name") || !strings.Contains(out, "id") {
		t.Fatalf("[%s] expected output to include projected headers name and id, got:\n%s", t.Name(), out)
	}

	// Should not include omitted header.
	if strings.Contains(out, "verified") {
		t.Fatalf("[%s] expected output not to include omitted header verified, got:\n%s", t.Name(), out)
	}

	// Basic sanity that row values are present.
	if !strings.Contains(out, "alice") || !strings.Contains(out, "bob") {
		t.Fatalf("[%s] expected output to include selected row values, got:\n%s", t.Name(), out)
	}
}
