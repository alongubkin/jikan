package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alongubkin/jikan/language"
)

func main() {
	parser := language.NewParser(strings.NewReader("every 100 sec of month"))

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
