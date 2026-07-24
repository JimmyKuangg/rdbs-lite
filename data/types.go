package data

type Database struct {
	Tables map[string]*Table
}

type Table struct {
	Name        string
	Schema      []Column
	Rows        []Row
	ColumnIndex map[string]int
}

type Column struct {
	Name string
	Type ColumnType
}

type Row struct {
	Values []any
}
