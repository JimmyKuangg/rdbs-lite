package commands

import (
	"errors"
	"rdbslite/data"
	"strings"
)

func Create(db *data.Database, cmd Command) (string, error) {
	if len(cmd.Args) < 2 {
		return "", errors.New("CREATE requires more arguments")
	}

	cmd.Args[0] = strings.ToUpper(cmd.Args[0])

	if cmd.Args[0] != "TABLE" {
		return "", errors.New("CREATE only supports TABLE currently")
	}

	schemaArgs := cmd.Args[2:]

	if len(schemaArgs)%2 != 0 || len(schemaArgs) == 0 {
		return "", errors.New("invalid schema")
	}

	columns := []data.Column{}

	for i := 0; i < len(schemaArgs); i += 2 {
		columns = append(columns, data.Column{
			Name: schemaArgs[i],
			Type: schemaArgs[i+1],
		})
	}

	err := db.CreateTable(cmd.Args[1], columns)
	if err != nil {
		return "", err
	}

	return "OK", nil
}
