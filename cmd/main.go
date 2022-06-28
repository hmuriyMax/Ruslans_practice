package main

import (
	"fmt"
	"time"
	"webScraper/pkg/data_parser"
)

func main() {
	//fmt.Print(getPriceOfOneItem("405000981"))
	start := time.Now()
	result, _ := data_parser.GetAllData("Shop")
	duration := time.Since(start)
	fmt.Println(duration)
	fmt.Println(duration.Nanoseconds())
	fmt.Print(result)
}
