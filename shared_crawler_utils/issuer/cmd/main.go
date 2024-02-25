package main

import (
	"encoding/csv"
	"os"
	"strings"

	utils "insiderviz.com/crawlers/shared_crawler_utils"
	issuer_utils "insiderviz.com/crawlers/shared_crawler_utils/issuer"
)

// Entry point to test functionality in the issuer package
func main() {
	file, err := os.Open("./data/tickers.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read in the csv
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	ciks := make([]string,len(data))
	for i,line := range data {
		// Each line has ticker\tcik
		cik := line[1]
		// left pads cik to 10 digits 1234567 -> 0001234567
		ciks[i] = strings.Repeat("0",10 - len(cik)) + cik
	}
	conn := utils.CreateMySQLConnection("root:root@tcp(127.0.0.1:3306)/insiderviz-crawler")
	issuer_utils.CrawlIssuersByCIK(conn,ciks,10)
}