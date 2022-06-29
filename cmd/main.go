package main

import (
	"log"
	"os"
	"time"
	"webScraper/pkg/data_parser"
)

func main() {
	logger := log.New(os.Stdout, "", 0)
	start := time.Now()
	result, err := data_parser.GetAllData(logger, "Shop", []string{"US", "GB", "CH"})
	duration := time.Since(start)
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
	data_parser.Print(logger, result)
	logger.Printf("Operation took %.0f seconds. Parsed %d items", duration.Seconds(), len(result))
}
