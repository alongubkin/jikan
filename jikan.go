package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	parser := NewParser(strings.NewReader("every 100 days"))

	every, err := parser.Parse()
	if err != nil {
		log.Println(err)
	}

	serialized, err := json.MarshalIndent(every, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	os.Stdout.Write(serialized)
}
