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
