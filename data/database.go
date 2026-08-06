package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func (db *Database) Save() error {
	path := filepath.Join(storagePath, dbFile)
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0o644,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	fmt.Println("Writing database to disk...")
	for _, table := range db.Tables {
		_, err = file.WriteString("TABLE " + table.Name + "\n")
		if err != nil {
			return err
		}

		_, err := file.WriteString("SCHEMA\n")
		if err != nil {
			return err
		}

		for _, column := range table.Schema {
			_, err = file.WriteString(string(column.Type) + " " + column.Name + "\n")
			if err != nil {
				return err
			}
		}

		_, err = file.WriteString("ROWS\n")
		if err != nil {
			return err
		}

		for _, row := range table.Rows {
			var rowStr strings.Builder

			for _, val := range row.Values {
				valStr := fmt.Sprintf("%v", val)
				_, err := rowStr.WriteString(valStr + " ")
				if err != nil {
					return err
				}
			}

			_, err = file.WriteString(rowStr.String() + "\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
