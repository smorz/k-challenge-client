package challenge

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

// NewCopier create a Copier instace
//
// source can be a random generator
func NewCopier(db *sql.DB, source Copyable) (*Copier, error) {
	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := txn.Prepare(pq.CopyIn(source.table(), source.fields()...))
	if err != nil {
		return nil, err
	}

	return &Copier{
		db:     db,
		txn:    txn,
		stmt:   stmt,
		source: source,
	}, nil
}

// Start generates statements concurrently and finally committs.
//
// workerCount is the number of workers in the workers pool;
// that can usually be set to the number of CPUs.
func (c *Copier) Start(workerCount int) (err error) {
	for i := 0; i < workerCount; i++ {
		c.wg.Add(1)
		go c.generateStatments()
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

// generateStatments
//
// as long as Copyable interface gives none-nil value,
// it receives values and make statements.
func (c *Copier) generateStatments() {
	defer c.wg.Done()
	for {
		v := c.source.values()
		if v == nil {
			return
		}
		_, err := c.stmt.Exec(v...)
		if err != nil {
			log.Fatal(err)
		}
	}
}
