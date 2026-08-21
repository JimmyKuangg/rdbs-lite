package repl

import (
	"errors"
	"strings"

	"rdbslite/commands"
	"rdbslite/data"
)

func ParseCommand(input string) (commands.Command, error) {
	fields := strings.Fields(input)

	if len(fields) == 0 {
		return commands.Command{}, errors.New("empty command")
	}

	command := strings.ToUpper(fields[0])

	return commands.Command{
		Name: command,
		Args: fields[1:],
	}, nil
}

func ExecuteCommand(db *data.Database, cmd commands.Command) (string, error) {
	switch cmd.Name {
	case "CREATE":
		resp, err := commands.Create(db, cmd)
		if err != nil {
			return "", err
		}
		if err := persistCommand(db, cmd); err != nil {
			return "", err
		}
		return resp, nil

	case "INSERT":
		resp, err := commands.Insert(db, cmd)
		if err != nil {
			return "", err
		}
		if err := persistCommand(db, cmd); err != nil {
			return "", err
		}
		return resp, nil

	case "PRINT":
		return commands.Print(db, cmd)

	case "SELECT":
		return commands.Select(db, cmd)

	case "SAVE":
		return commands.Save(db)

	default:
		return "", errors.New("unknown command")
	}
}

func persistCommand(db *data.Database, cmd commands.Command) error {
	if err := AppendAOF(cmd.ToString()); err != nil {
		return err
	}
	return checkpointIfNeeded(db)
}
