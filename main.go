package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// db := data.NewDatabase()
	reader := bufio.NewScanner(os.Stdin)

	fmt.Println("RDBSLite started")

	for reader.Scan() {
		text := strings.TrimSpace(reader.Text())

		fmt.Println(text)
	}

	if err := reader.Err(); err != nil {
		log.Fatalf("uh oh! %v", err)
	}
}
