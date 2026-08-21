package data

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (db *Database) Save() error {
	tmpFile, err := os.CreateTemp(storagePath, "storage-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	success := false
	defer func() {
		_ = tmpFile.Close()
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()

	fmt.Println("Writing database to disk...")
	for _, table := range db.Tables {
		_, err = tmpFile.WriteString("TABLE\n" + table.Name + "\n")
		if err != nil {
			return err
		}

		_, err = tmpFile.WriteString("SCHEMA\n")
		if err != nil {
			return err
		}

		for _, column := range table.Schema {
			_, err = tmpFile.WriteString(string(column.Type) + " " + column.Name + "\n")
			if err != nil {
				return err
			}
		}

		_, err = tmpFile.WriteString("ROWS\n")
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

			_, err = tmpFile.WriteString(rowStr.String() + "\n")
			if err != nil {
				return err
			}
		}
	}

	if err := tmpFile.Sync(); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	path := filepath.Join(storagePath, dbFile)
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	success = true
	return nil
}

func (db *Database) Load() error {
	dbFilePath := filepath.Join(storagePath, dbFile)
	file, err := os.OpenFile(
		dbFilePath,
		os.O_RDONLY,
		0o644,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	var currentTable *Table
	section := ""
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch line {
		case "TABLE":
			section = "table"
		case "SCHEMA":
			section = "schema"
		case "ROWS":
			section = "rows"

		default:
			switch section {
			case "table":
				currentTable = &Table{
					Name:        line,
					ColumnIndex: make(map[string]int),
				}
				db.Tables[strings.ToLower(line)] = currentTable

			case "schema":
				parts := strings.Fields(line)
				colType, err := ParseColumnType(parts[0])
				if err != nil {
					return err
				}

				col := Column{Name: parts[1], Type: colType}
				currentTable.Schema = append(currentTable.Schema, col)
				currentTable.ColumnIndex[strings.ToLower(col.Name)] = len(currentTable.Schema) - 1

			case "rows":
				fields := strings.Fields(line)
				values := make([]any, len(fields))

				for i, f := range fields {
					v, err := ParseValue(f, currentTable.Schema[i].Type)
					if err != nil {
						return err
					}

					values[i] = v
				}
				currentTable.Rows = append(currentTable.Rows, Row{Values: values})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
