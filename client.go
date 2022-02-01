package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
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
}
