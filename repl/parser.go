package repl

import (
	"errors"
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
