package main

import (
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
	live "insiderviz.com/crawlers/live_crawler"
	local "insiderviz.com/crawlers/local_crawler"
	utils "insiderviz.com/crawlers/utils"
	stock_data_utils "insiderviz.com/crawlers/utils/issuer/stock_data"
)

// This file can run any crawler, should mainly be used for local testing as crawlers should be deployed indiviually
func main() {
	godotenv.Load()
	// Creates connection to local db called insiderviz-crawler2 with username and password "root"
	conn := utils.CreateMySQLConnection("root:root@tcp(127.0.0.1:3306)/insiderviz-crawler2") 
	defer conn.Close()
	// cik := "0001378992"
	// issuer_utils.CrawlIssuerJSON(cik)
	test_live(conn)
	// test_local(conn)
	// crawl_all_stock_data(conn)
}

func test_live(conn *sql.DB) {
	live.LiveCrawl(conn)
}

func test_local(conn *sql.DB) {
	folder := "./submissions/2024-03-04"

	offset := 0 // start at the first company
	stride := 1 // only go one at a time
	start := "2024-01-26" //start 
	end := "3000-01-01" // end, arbitrary date in future means get everything
	local.RunIssuerCrawl(folder,conn)
	local.RunFormCrawl(folder,start,end,offset,stride,conn)
}

// Helpful function that will populate the stock_data table using the EOD Historical Data API and the tickers table in the conn mysql db
func crawl_all_stock_data(conn * sql.DB) {

	// Reads all tickers currently in DB
	tickers := make([]string,0)
	get_ticker_sql := `
	SELECT ticker from ticker
	`
	result,err := conn.Query(get_ticker_sql)
	if err != nil {
		panic(err)
	}
	for result.Next() {
		var ticker string
		result.Scan(&ticker)
		tickers = append(tickers,ticker)
	}
	fmt.Printf("Found %d tickers\n",len(tickers))
	// Saves stock data for all tickers
	stock_data_utils.CrawlStockData(conn,tickers,"0001-01-01",20)
}