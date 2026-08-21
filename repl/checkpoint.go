package repl

import (
	"os"

	"rdbslite/data"
)

const aofCheckpointThreshold = 256

func checkpointIfNeeded(db *data.Database) error {
	aofPath := data.AOFPath()
	info, err := os.Stat(aofPath)
	if err != nil {
		return err
	}

	if info.Size() >= aofCheckpointThreshold {
		if err := db.Save(); err != nil {
			return err
		}
		if err := os.Truncate(aofPath, 0); err != nil {
			return err
		}
	}

	// Note for self -
	// There is an atomicity / durability  issue here between the save and the truncate
	// If we crash between save and truncate, we will have an issue with replaying commands that can error out
	// For instance, we create a table -> The write command gets written to AOF -> we run to save -> crash
	// Our write to disk will contain our newest table, but the AOF will not be flushed, leading to us trying to replay the cmd that creates the table
	// This can cause us to error out in the middle of our replay because the table already exists
	// Issue to look at or think about down the line

	return nil
}
