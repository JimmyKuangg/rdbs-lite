package data

import "errors"

func (db *Database) CreateTable(name string) error {
	t := &Table{
		Name: name,
	}

	if _, exists := db.Tables[name]; exists {
		return errors.New("table already exists")
	}

	db.Tables[t.Name] = t

	return nil
}
