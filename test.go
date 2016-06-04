package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alongubkin/jikan/parser"
)

func main() {
	parser := parser.NewParser(strings.NewReader("every minute"))

	every, err := parser.ParseIntervalsSchedule()
	if err != nil {
		log.Println(err)
	}

	serialized, err := json.MarshalIndent(every, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	os.Stdout.Write(serialized)
}
