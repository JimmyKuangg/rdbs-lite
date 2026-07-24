package commands

import (
	"rdbslite/data"
	"testing"
)

func TestInsert(t *testing.T) {
	dbCreateCmd := Command{
		Name: "CREATE",
		Args: []string{"TABLE", "users", "id", "int"},
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
