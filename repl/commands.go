package repl

import (
	"errors"
	"rdbslite/data"
	"strings"
)

func ParseCommand(input string) (Command, error) {
	fields := strings.Fields(input)

	if len(fields) == 0 {
		return Command{}, errors.New("empty command")
	}

	command := strings.ToUpper(fields[0])

	return Command{
		Name: command,
		Args: fields[1:],
	}, nil
}

func ExecuteCommand(db *data.Database, cmd Command) (string, error) {
	switch cmd.Name {

	case "CREATE":
		if len(cmd.Args) < 2 {
			return "", errors.New("CREATE requires more arguments")
		}

		cmd.Args[0] = strings.ToUpper(cmd.Args[0])

		if cmd.Args[0] != "TABLE" {
			return "", errors.New("CREATE only supports TABLE currently")
		}

		err := db.CreateTable(cmd.Args[1])
		if err != nil {
			return "", err
		}

		return "OK", nil

	default:
		return "", errors.New("unknown command")
	}
}
