package commands

import (
	"errors"
	"fmt"
	"rdbslite/data"
	"strings"
)

func Create(db *data.Database, cmd Command) (string, error) {
	if len(cmd.Args) < 2 {
		return "", errors.New("CREATE requires more arguments")
	}

	createSubj := strings.ToUpper(cmd.Args[0])

	if createSubj != "TABLE" {
		return "", errors.New("CREATE only supports TABLE currently")
	}

	schemaArgs := cmd.Args[2:]

	if len(schemaArgs)%2 != 0 || len(schemaArgs) == 0 {
		return "", errors.New("invalid schema")
	}

	columns := []data.Column{}
	seen := make(map[string]bool)

	for i := 0; i < len(schemaArgs); i += 2 {
		columnName := strings.TrimSpace(schemaArgs[i])
		if columnName == "" {
			return "", errors.New("column name cannot be empty")
		}

		if data.IsReservedIdentifier(columnName) {
			return "", fmt.Errorf("can not use name %q as a column name: reserved keyword", columnName)
		}

		key := strings.ToLower(columnName)
		if _, exists := seen[key]; exists {
			return "", fmt.Errorf("duplicate column name: %q", columnName)
		}
		seen[key] = true

		columnType, err := data.ParseColumnType(schemaArgs[i+1])
		if err != nil {
			return "", err
		}

		columns = append(columns, data.Column{
			Name: columnName,
			Type: columnType,
		})
	}

	err := db.CreateTable(cmd.Args[1], columns)
	if err != nil {
		return "", err
	}

	return "OK", nil
}
