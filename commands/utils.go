package commands

import "strings"

func (cmd Command) ToString() string {
	return strings.Join(append([]string{cmd.Name}, cmd.Args...), " ")
}
