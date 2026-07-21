package repl

import "rdbslite/data"

type Command struct {
	Name string
	Args []string
}

type TableDef struct {
	Name   string
	Schema []data.Column
}
