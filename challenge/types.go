package challenge

import "database/sql"

type Copyable interface {
	table() string
	fields() []interface{}
	values() ([]interface{}, bool)
}

type Generator struct {
	db *sql.DB
}
