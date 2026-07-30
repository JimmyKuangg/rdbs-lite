package data

import (
	"errors"
	"fmt"
	"strings"
)

func (db *Database) CreateTable(name string, schema []Column) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.New("table name cannot be empty")
	}

	key := strings.ToLower(trimmed)
	if _, exists := db.Tables[key]; exists {
		return errors.New("table already exists")
	}

	if IsReservedIdentifier(trimmed) {
		return fmt.Errorf("can not use name %v as a table name: reserved keyword", trimmed)
	}

	t := &Table{
		Name:        trimmed,
		Schema:      schema,
		ColumnIndex: make(map[string]int),
	}

	for i, col := range schema {
		t.ColumnIndex[strings.ToLower(col.Name)] = i
	}

	db.Tables[key] = t
	return nil
}

func (db *Database) Insert(tableName string, row []any) error {
	table := db.Tables[tableName]
	stored := append([]any(nil), row...)
	table.Rows = append(table.Rows, Row{Values: stored})
	return nil
}

func (db *Database) Select(table Table, selectedCols []string, where []string) ([][]any, error) {
	results := make([][]any, 0)

	if len(where) == 0 {
		for _, row := range table.Rows {
			projected := make([]any, 0, len(selectedCols))
			for _, col := range selectedCols {
				idx := table.ColumnIndex[strings.ToLower(col)]
				projected = append(projected, row.Values[idx])
			}
			results = append(results, projected)
		}
		return results, nil
	}

	whereTargetIdx := table.ColumnIndex[strings.ToLower(where[0])]
	whereTargetType := table.Schema[whereTargetIdx].Type

	target, err := ParseValue(where[2], whereTargetType)
	if err != nil {
		return nil, fmt.Errorf("error during SELECT: %w", err)
	}

	for _, row := range table.Rows {
		rowTarget := row.Values[whereTargetIdx]
		if compareValues(rowTarget, where[1], target) {
			projected := make([]any, 0, len(selectedCols))
			for _, col := range selectedCols {
				idx := table.ColumnIndex[strings.ToLower(col)]
				projected = append(projected, row.Values[idx])
			}
			results = append(results, projected)
		}
	}

	return results, nil
}
