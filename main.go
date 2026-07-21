package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"rdbslite/data"
	"rdbslite/repl"
	"strings"
)

func main() {
	db := data.NewDatabase()
	reader := bufio.NewScanner(os.Stdin)

	fmt.Println("RDBSLite started")

	for reader.Scan() {
		text := strings.TrimSpace(reader.Text())

		command, err := repl.ParseCommand(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err = repl.ExecuteCommand(&db, command); err != nil {
			fmt.Println(err)
		}
	}

	if err := reader.Err(); err != nil {
		log.Fatalf("uh oh! %v", err)
	}
}
