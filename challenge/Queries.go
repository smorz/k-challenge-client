package challenge

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

func NewCopier(db *sql.DB, tabale Copyable) *Copier {
	txn, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}
	stmt, err := txn.Prepare(pq.CopyIn(tabale.table(), tabale.fields()...))
	if err != nil {
		log.Fatal(err)
	}
	return &Copier{
		db:   db,
		txn:  txn,
		stmt: stmt,
	}
}

func (c *Copier) Start() error {

}
