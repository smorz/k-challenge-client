package challenge

import (
	"database/sql"
	"time"
)

type Copyable interface {
	table() string
	fields() []string
	values() ([]interface{}, bool)
}

type Trade struct {
	firstDay     time.Time
	days         chan TradableDay
	generateDone bool
}

type TradableDay struct {
	instrumentID int
	deyOffset    int
}

type Generator struct {
	db *sql.DB
}
