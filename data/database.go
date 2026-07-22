package data

import (
	"errors"
	"strings"
)

func (db *Database) CreateTable(name string, schema []Column) error {
	if _, exists := db.Tables[strings.ToLower(name)]; exists {
		return errors.New("table already exists")
	}

	t := &Table{
		Name:   name,
		Schema: schema,
	}

	db.Tables[name] = t

	return nil
}
