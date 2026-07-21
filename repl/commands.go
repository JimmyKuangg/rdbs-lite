package repl

import (
	"errors"
	"rdbslite/commands"
	"rdbslite/data"
	"strings"
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
		return commands.Create(db, cmd)

	default:
		return "", errors.New("unknown command")
	}
}
