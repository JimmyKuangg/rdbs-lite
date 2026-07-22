package data

import (
	"errors"
)

func (db *Database) CreateTable(name string, schema []Column) error {
	if _, exists := db.Tables[name]; exists {
		return errors.New("table already exists")
	}

	t := &Table{
		Name:   name,
		Schema: schema,
	}

	db.Tables[name] = t

	return nil
}
