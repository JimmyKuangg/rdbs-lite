package commands

import (
	"fmt"

	"rdbslite/data"
)

func Save(db *data.Database) (string, error) {
	err := db.Save()
	if err != nil {
		return "", fmt.Errorf("error during db save: %w", err)
	}

	return "OK", nil
}
