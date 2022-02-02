package challenge

import (
	"database/sql"
	"sync"
	"time"
)

// Copyable is any type that can be copy by postgres.
type Copyable interface {
	table() string
	fields() []string
	values() []interface{}
}

// TradeGenerator is a Copyable type that generates randome recoreds.
type TradeGenerator struct {
	recordsCount int
	generated    int
	firstDay     time.Time
	days         chan instrumentDay
	mu           sync.Mutex
}

// instrumentDay shows what instrument traded and in how many days after the start day?
type instrumentDay struct {
	instrumentID int
	deyOffset    int
}

// Copier by using a copyable type generates copy statements concurrently
// and finally commit.
type Copier struct {
	db     *sql.DB
	txn    *sql.Tx
	stmt   *sql.Stmt
	source Copyable
	wg     sync.WaitGroup
}
