package data

func NewDatabase() Database {
	return Database{
		Tables: make(map[string]*Table),
	}
}
