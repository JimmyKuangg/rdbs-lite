package commands

import (
	"rdbslite/data"
	"testing"
)

func TestInsert(t *testing.T) {
	dbCreateCmd := Command{
		Name: "CREATE",
		Args: []string{"TABLE", "users", "id", "int", "name", "text", "verified", "bool"},
	}

	tests := []TestCases{
		{
			name: "inserts data into a table",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "users", "id", "1"},
			},
			wantErr: false,
		},
		{
			name: "rejects with not enough arguments",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "users"},
			},
			wantErr: true,
		},
		{
			name: "rejects with invalid table name",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "definitelyarealtable", "id", "1"},
			},
			wantErr: true,
		},
		{
			name: "rejects multiple inserts into the same column",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "users", "id", "1", "id", "2"},
			},
			wantErr: true,
		},
		{
			name: "rejects with invalid column value pairs",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "users", "id"},
			},
			wantErr: true,
		},
		{
			name: "rejects with invalid column name",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "users", "definitelyarealcolumn", "2"},
			},
			wantErr: true,
		},
		{
			name: "rejects with invalid column column typing for int",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "users", "id", "jimbolina"},
			},
			wantErr: true,
		},
		{
			name: "rejects with invalid column column typing for bool",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"INTO", "users", "verified", "sometimes"},
			},
			wantErr: true,
		},
		{
			name: "works with case insensitvity",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"into", "USERS", "ID", "2"},
			},
			wantErr: false,
		},
		{
			name: "works with unordered columns",
			cmd: Command{
				Name: "INSERT",
				Args: []string{"into", "USERS", "name", "Bob", "id", "1"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := data.NewDatabase()
			_, err := Create(&db, dbCreateCmd)
			if err != nil {
				t.Fatalf("error creating db")
			}

			_, err = Insert(&db, tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsert_RowOrder(t *testing.T) {
	db := data.NewDatabase()
	dbCreateCmd := Command{
		Name: "CREATE",
		Args: []string{"TABLE", "users", "id", "int", "name", "text", "verified", "bool"},
	}

	_, err := Create(&db, dbCreateCmd)
	if err != nil {
		t.Fatalf("[%s] error creating db: %v", t.Name(), err)
	}

	insertCmd := Command{
		Name: "INSERT",
		Args: []string{"INTO", "users", "name", "realname", "verified", "true", "id", "2"},
	}
	_, err = Insert(&db, insertCmd)
	if err != nil {
		t.Fatalf("[%s] insert failed: %v", t.Name(), err)
	}

	table := db.Tables["users"]
	if len(table.Rows) != 1 {
		t.Fatalf("[%s] expected 1 row, got %d", t.Name(), len(table.Rows))
	}

	got := table.Rows[0].Values
	if len(got) != 3 {
		t.Fatalf("[%s] expected 3 values, got %d", t.Name(), len(got))
	}

	if got[0] != 2 {
		t.Fatalf("[%s] expected id at index 0 to be 2, got %v", t.Name(), got[0])
	}
	if got[1] != "realname" {
		t.Fatalf("[%s] expected name at index 1 to be realname, got %v", t.Name(), got[1])
	}
	if got[2] != true {
		t.Fatalf("[%s] expected verified at index 2 to be true, got %v", t.Name(), got[2])
	}
}
