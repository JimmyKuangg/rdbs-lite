package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"rdbslite/data"
	"rdbslite/repl"
)

func main() {
	if err := data.Init(); err != nil {
		log.Fatal(err)
	}

	db := data.NewDatabase()
	if err := db.Load(); err != nil {
		log.Fatalf("error loading database: %v", err)
	}
	reader := bufio.NewScanner(os.Stdin)

	fmt.Println("RDBSLite started")

	for reader.Scan() {
		text := strings.TrimSpace(reader.Text())

		command, err := repl.ParseCommand(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		resp, err := repl.ExecuteCommand(&db, command)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(resp)
	}

	if err := reader.Err(); err != nil {
		log.Fatalf("uh oh! %v", err)
	}
}
