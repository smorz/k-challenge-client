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
	startMoment := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
	rand.Seed(time.Now().UnixNano())

	firstDay := time.Now().AddDate(-1, 0, 0)
	tg, _ := challenge.NewTradeGenerator(firstDay, count)

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	c, err := challenge.NewCopier(db, tg)
	if err != nil {
		log.Fatal(err)
	}

	cores := runtime.GOMAXPROCS(-1)

	if err := c.Start(cores); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("elapsed time: %v\n", time.Since(startMoment))

}
