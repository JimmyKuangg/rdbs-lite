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

func ExecuteCommand(db *data.Database, cmd Command) error {
	switch cmd.Name {

	case "CREATE":
		if len(cmd.Args) < 2 {
			return errors.New("CREATE requires more arguments")
		}

	default:
		return errors.New("unknown command")
	}

	return nil
}
