package main

import (
	"log"
	"os"
	"time"
	"webScraper/pkg/data_parser"
)

func main() {
	start := time.Now()
	logger := log.New(os.Stdout, "", 0)
	result, _ := data_parser.GetAllData(logger, "Shop")
	duration := time.Since(start)
	logger.Println(result)
	logger.Printf("Operation took %.0f seconds. Parsed %d items", duration.Seconds(), len(result))
}
