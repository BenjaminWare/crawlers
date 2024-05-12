package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	issuer_crawler "insiderviz.com/crawlers/issuer_crawler"
	live "insiderviz.com/crawlers/live_crawler"
	local "insiderviz.com/crawlers/local_crawler"
	stock_day_crawler "insiderviz.com/crawlers/stock_day_crawler"
	utils "insiderviz.com/crawlers/utils"
	stock_data_utils "insiderviz.com/crawlers/utils/issuer/stock_data"
)

// This file can run any crawler, should mainly be used for local testing as crawlers should be deployed indiviually
func main() {
	godotenv.Load()
	// Creates connection to local db called insiderviz-crawler2 with username and password "root"
	conn := utils.CreateMySQLConnection(os.Getenv("CONNECTION_STRING"))
	defer conn.Close()
	// cik := "0001378992"
	// issuer_utils.CrawlIssuerJSON(cik)
	// test_live(conn)
	// test_local(conn)
	crawl_all_stock_data(conn)
	// test_stock_day_crawler(conn)
}

func test_stock_day_crawler(conn *sql.DB) {
	stock_day_crawler.CrawlTodaysStockDay(conn, 20)
}
func test_issuer_crawler(conn *sql.DB) {
	issuer_crawler.CrawlIssuers(conn, 20)
}

func test_live(conn *sql.DB) {
	live.LiveCrawl(conn)
}

func test_local(conn *sql.DB) {
	folder := "./submissions"

	offset := 0                        // start at the first company
	stride := 1                        // only go one at a time
	start := "2024-03-04"              //start
	end := "3000-01-01"                // end, arbitrary date in future means get everything
	local.RunIssuerCrawl(folder, conn) // Find issuers we don't have already
	local.RunFormCrawl(folder, start, end, offset, stride, conn)
}

// Helpful function that will populate the stock_data table using the EOD Historical Data API and the tickers table in the conn mysql db
func crawl_all_stock_data(conn *sql.DB) {

	// Reads all tickers currently in DB
	tickers := make([]string, 0)
	get_ticker_sql := `
	SELECT ticker from ticker
	`
	result, err := conn.Query(get_ticker_sql)
	if err != nil {
		panic(err)
	}
	for result.Next() {
		var ticker string
		result.Scan(&ticker)
		tickers = append(tickers, ticker)
	}
	fmt.Printf("Found %d tickers\n", len(tickers))
	// Saves stock data for all tickers
	stock_data_utils.CrawlStockData(conn, tickers, "2024-03-06", 20)
}
