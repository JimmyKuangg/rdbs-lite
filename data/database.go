package data

import "errors"

func (db *Database) CreateTable(name string) error {
	if _, exists := db.Tables[name]; exists {
		return errors.New("table already exists")
	}

	t := &Table{
		Name: name,
	}

	db.Tables[t.Name] = t

	return nil
}
