package challenge

import (
	"database/sql"
	"sync"
	"time"
)

type Copyable interface {
	table() string
	fields() []string
	values() []interface{}
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

type Copier struct {
	db     *sql.DB
	txn    *sql.Tx
	stmt   *sql.Stmt
	source Copyable
	wg     sync.WaitGroup
}
