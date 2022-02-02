package challenge

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

func NewCopier(db *sql.DB, table Copyable) (*Copier, error) {
	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := txn.Prepare(pq.CopyIn(table.table(), table.fields()...))
	if err != nil {
		return nil, err
	}

	return &Copier{
		db:    db,
		txn:   txn,
		stmt:  stmt,
		table: table,
	}, nil
}

func (c *Copier) Start(cpuCount int) (err error) {
	for i := 0; i < cpuCount; i++ {
		c.wg.Add(1)
		go c.GenerateStatments()
	}
	c.wg.Wait()
	_, err = c.stmt.Exec()
	if err != nil {
		return
	}
	if err = c.stmt.Close(); err != nil {
		log.Printf("error during stmt.Close(): %s\n", err)
		return
	}
	if err = c.txn.Commit(); err != nil {
		log.Printf("could not commit transaction: %s\n", err)
		return
	}
	return
}

func (c *Copier) GenerateStatments() {
	defer c.wg.Done()
	for {
		v := c.table.values()
		if v == nil {
			return
		}
		_, err := c.stmt.Exec(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}