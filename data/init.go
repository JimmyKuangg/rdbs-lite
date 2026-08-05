package data

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	storagePath = ".rdbslite"
	dbFile      = "storage.db"
)

func NewDatabase() Database {
	return Database{
		Tables: make(map[string]*Table),
	}
}

func Init() error {
	err := os.MkdirAll(storagePath, 0o755)
	if err != nil {
		return fmt.Errorf("error creating the rdbslite directory: %w", err)
	}

	storageFilePath := filepath.Join(storagePath, dbFile)
	file, err := os.OpenFile(
		storageFilePath,
		os.O_CREATE|os.O_RDWR,
		0o644,
	)
	if err != nil {
		return fmt.Errorf("error during init: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	return nil
}
