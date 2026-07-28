package commands

import (
	"errors"
	"rdbslite/data"
	"strings"
)

func Select(db *data.Database, cmd Command) (string, error) {
	if len(cmd.Args) != 3 {
		return "", errors.New("SELECT requires more arguments")
	}

	if cmd.Args[0] != "*" {
		return "", errors.New("currently only supports all (*)")
	}

	var fromKeywordIdx int

	for i := range cmd.Args {
		if strings.ToUpper(cmd.Args[i]) == "FROM" {
			fromKeywordIdx = i
		}
	}

	if fromKeywordIdx == 0 {
		return "", errors.New("keyword FROM required in SELECT statement")
	}

	tableName := strings.ToLower(cmd.Args[2])
	if _, exists := db.Tables[tableName]; !exists {
		return "", errors.New("table does not exist")
	}

	return renderTable(db.Tables[tableName]), nil
}
