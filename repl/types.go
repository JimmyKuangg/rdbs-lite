package repl

import "rdbslite/data"

type TableDef struct {
	Name   string
	Schema []data.Column
}
