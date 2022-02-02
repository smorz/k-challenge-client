/* A solution to a code challenge that wanted to insert
a number of random records with the highest efficiency
in the trading table. A search revealed that in
Postgres, copy is the best option.*/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/smorz/k-challenge-client/challenge"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "123456"
	DB_NAME     = "ktest"
)

func main() {
	// just for studying program efficiency
	startMoment := time.Now()

	// setup log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check the accuracy of the input
	if len(os.Args) < 2 {
		log.Fatal("The command must have a positive integer input next to it")
	}
	if len(os.Args) > 2 {
		log.Fatal("The command must have only one input. no more")
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil || count <= 0 {
		log.Fatal("The input must be a positive integer")
	}

	// used in TradeGenerator
	rand.Seed(time.Now().UnixNano())

	firstDay := time.Now().AddDate(-1, 0, 0)

	// Creating an instance of generator
	tg, _ := challenge.NewTradeGenerator(firstDay, count)

	// connecting to database
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("connected to postgresql")

	// setup database
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	// Creating an instance of Copier
	c, err := challenge.NewCopier(db, tg)
	if err != nil {
		log.Fatal(err)
	}

	// Obtain the number of CPUs
	cores := runtime.GOMAXPROCS(-1)

	// Start Generating and inserting
	// Each row is generate just before the generation of the copy statement.
	if err := c.Start(cores); err != nil {
		log.Fatal(err)
	}

	// showing Result
	fmt.Printf("elapsed time: %v\n", time.Since(startMoment))

}
