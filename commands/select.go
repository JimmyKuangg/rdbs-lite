package commands

import (
	"errors"
	"fmt"
	"rdbslite/data"
	"strings"
)

func Select(db *data.Database, cmd Command) (string, error) {
	if len(cmd.Args) < 3 {
		return "", errors.New("SELECT requires more arguments")
	}

	fromKeywordIdx := -1

	for i := range cmd.Args {
		if strings.ToUpper(cmd.Args[i]) == "FROM" {
			fromKeywordIdx = i
			break
		}
	}

	if fromKeywordIdx == -1 {
		return "", errors.New("keyword FROM required in SELECT statement")
	}

	if fromKeywordIdx == 0 {
		return "", errors.New("SELECT requires at least one column before FROM")
	}

	if fromKeywordIdx+1 >= len(cmd.Args) {
		return "", errors.New("missing table name after FROM")
	}

	if fromKeywordIdx+2 != len(cmd.Args) {
		return "", errors.New("unexpected tokens after table name")
	}

	tableName := strings.ToLower(cmd.Args[fromKeywordIdx+1])
	table := db.Tables[tableName]
	if table == nil {
		return "", errors.New("table does not exist")
	}

	columns := cmd.Args[:fromKeywordIdx]

	resolvedCols, err := resolveProjection(table, columns)
	if err != nil {
		return "", err
	}

	return renderTable(table, resolvedCols), nil
}

func resolveProjection(table *data.Table, columns []string) ([]string, error) {
	if len(columns) == 1 && columns[0] == "*" {
		return []string{"*"}, nil
	}

	resolved := make([]string, len(columns))
	for i, name := range columns {
		col := strings.ToLower(name)
		if _, exists := table.ColumnIndex[col]; !exists {
			return nil, fmt.Errorf("column name does not exist in %s table", table.Name)
		}
		resolved[i] = col
	}

	return resolved, nil
}
