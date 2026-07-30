package commands

import (
	"errors"
	"fmt"
	"strings"

	"rdbslite/data"
)

func Select(db *data.Database, cmd Command) (string, error) {
	if len(cmd.Args) < 3 {
		return "", errors.New("SELECT requires more arguments")
	}

	fromKeywordIdx := parseClause(cmd, "FROM")
	whereKeywordIdx := parseClause(cmd, "WHERE")

	if fromKeywordIdx == -1 {
		return "", errors.New("keyword FROM required in SELECT statement")
	}

	if fromKeywordIdx == 0 {
		return "", errors.New("SELECT requires at least one column before FROM")
	}

	if fromKeywordIdx+1 >= len(cmd.Args) {
		return "", errors.New("missing table name after FROM")
	}

	if whereKeywordIdx == -1 {
		if fromKeywordIdx+2 != len(cmd.Args) {
			return "", errors.New("unexpected tokens after table name")
		}
	} else {
		if fromKeywordIdx+2 != whereKeywordIdx {
			return "", errors.New("unexpected tokens before WHERE")
		}
		if whereKeywordIdx+1 >= len(cmd.Args) {
			return "", errors.New("missing condition after WHERE")
		}
	}

	tableName := strings.ToLower(cmd.Args[fromKeywordIdx+1])
	table := db.Tables[tableName]
	if table == nil {
		return "", errors.New("table does not exist")
	}

	columns := cmd.Args[:fromKeywordIdx]
	var whereArgs []string

	resolvedCols, err := resolveColumns(table, columns)
	if err != nil {
		return "", err
	}

	if whereKeywordIdx != -1 {
		whereWords := cmd.Args[whereKeywordIdx+1:]
		err = resolveWhereTokens(table, whereWords)
		if err != nil {
			return "", err
		}

		whereArgs = whereWords
	}

	rows, err := db.Select(*table, resolvedCols, whereArgs)
	if err != nil {
		return "", err
	}

	return renderRows(resolvedCols, rows), nil
}

func resolveColumns(table *data.Table, columns []string) ([]string, error) {
	if len(columns) == 1 && columns[0] == "*" {
		resolved := make([]string, len(table.Schema))
		for i, col := range table.Schema {
			resolved[i] = strings.ToLower(col.Name)
		}
		return resolved, nil
	}

	resolved := make([]string, len(columns))
	seen := make(map[string]bool)

	for i, name := range columns {
		col := strings.ToLower(name)
		if seen[col] {
			return nil, fmt.Errorf("can not use duplicate column names in SELECT")
		}
		seen[col] = true

		if _, exists := table.ColumnIndex[col]; !exists {
			return nil, fmt.Errorf("column name does not exist in %s table", table.Name)
		}
		resolved[i] = col
	}

	return resolved, nil
}

func resolveWhereTokens(table *data.Table, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("not enough args after WHERE clause")
	}

	colName := strings.ToLower(args[0])
	if _, exists := table.ColumnIndex[colName]; !exists {
		return fmt.Errorf("column name %s does not exist in table %s", args[0], table.Name)
	}

	op := args[1]
	switch op {
	case "=", "<", ">", "<=", ">=":
		// valid operator
	default:
		return fmt.Errorf("unsupported operator %s", op)
	}

	return nil
}

func parseClause(cmd Command, clause string) int {
	idx := -1

	for i, arg := range cmd.Args {
		if strings.EqualFold(arg, clause) {
			idx = i
			break
		}
	}

	return idx
}
