package challenge

import (
	"database/sql"
	"sync"
	"time"
)

type Copyable interface {
	table() string
	fields() []string
	values() ([]interface{}, bool)
	done()
	count() int
}

type TradeGenerator struct {
	recordsCount int
	generated    int
	firstDay     time.Time
	days         chan TradableDay
	mu           sync.Mutex
}

type TradableDay struct {
	instrumentID int
	deyOffset    int
}

type Generator struct {
	db *sql.DB
}
