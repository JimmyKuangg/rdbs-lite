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

	if len(cmd.Args)%2 != 0 {
		return "", errors.New("invalid schema")
	}

	err := db.CreateTable(cmd.Args[1])
	if err != nil {
		return "", err
	}

	return "OK", nil
}
