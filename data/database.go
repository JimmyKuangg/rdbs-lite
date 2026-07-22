package data

import (
	"errors"
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

	t := &Table{
		Name:   trimmed,
		Schema: schema,
	}

	db.Tables[key] = t
	return nil
}
