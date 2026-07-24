package commands

import (
	"errors"
	"fmt"
	"rdbslite/data"
	"strings"
)

func Insert(db *data.Database, cmd Command) (string, error) {
	if len(cmd.Args) < 2 {
		return "", errors.New("INSERT requires more arguments")
	}

	intoKeyword := strings.ToUpper(cmd.Args[0])
	if intoKeyword != "INTO" {
		return "", errors.New("INTO keyword required after INSERT")
	}

	tableName := strings.ToLower(cmd.Args[1])
	if _, exists := db.Tables[tableName]; !exists {
		return "", errors.New("table does not exist")
	}

	insertArgs := cmd.Args[2:]
	if len(insertArgs) == 0 || len(insertArgs)%2 == 1 {
		return "", fmt.Errorf("invalid INSERT syntax: expected pairs of <column> <value> after table name, got %v arguments", len(insertArgs))
	}

	table := db.Tables[tableName]
	seenCols := make(map[string]bool)
	newRow := make([]any, len(table.Schema))

	for i := 0; i < len(insertArgs); i += 2 {
		col := strings.ToLower(insertArgs[i])
		val := insertArgs[i+1]

		if seenCols[col] {
			return "", fmt.Errorf("invalid INSERT syntax: can not use the same column: '%v'", col)
		}
		seenCols[col] = true

		index, exists := table.ColumnIndex[col]
		if !exists {
			return "", fmt.Errorf("column does not exist: '%v'", col)
		}

		colType := table.Schema[index].Type
		res, err := data.ParseValue(val, colType)
		if err != nil {
			return "", fmt.Errorf("invalid typing for INSERT: column: '%v', value: '%v'", col, val)
		}

		newRow[index] = res
	}

	db.Insert(tableName, newRow)
	return "OK", nil
}
