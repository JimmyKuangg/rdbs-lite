package repl

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"rdbslite/commands"
	"rdbslite/data"
)

func Replay(db *data.Database) (err error) {
	aofPath := data.AOFPath()
	file, err := os.OpenFile(
		aofPath,
		os.O_APPEND|os.O_RDONLY,
		0o644,
	)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		cmd, err := ParseCommand(line)
		if err != nil {
			return err
		}

		if err := replayCommand(db, cmd); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func replayCommand(db *data.Database, cmd commands.Command) error {
	switch cmd.Name {
	case "CREATE":
		_, err := commands.Create(db, cmd)
		return err

	case "INSERT":
		_, err := commands.Insert(db, cmd)
		return err

	default:
		return errors.New("unknown command")
	}
}
