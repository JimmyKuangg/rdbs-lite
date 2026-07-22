package data

type Database struct {
	Tables map[string]*Table
}

type Table struct {
	Name   string
	Schema []Column
}

type Column struct {
	Name string
	Type ColumnType
}
