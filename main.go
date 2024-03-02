package main

import (
	"database/sql"

	live "insiderviz.com/crawlers/live_crawler"
	local "insiderviz.com/crawlers/local_crawler"
	utils "insiderviz.com/crawlers/shared_crawler_utils"
)

// This file can run any crawler, should mainly be used for local testing as crawlers should be deployed indiviually
func main() {
	// Creates connection to local db called insiderviz-crawler2 with username and password "root"
	conn := utils.CreateMySQLConnection("root:root@tcp(127.0.0.1:3306)/insiderviz-crawler") 
	defer conn.Close()
	// cik := "0001378992"
	// issuer_utils.CrawlIssuerJSON(cik)
	// test_live(conn)
	test_local(conn)
}

func test_live(conn *sql.DB) {
	live.RunLiveCrawler(conn)
}

func test_local(conn *sql.DB) {
	folder := "./submissions/2024-01-27"

	offset := 4998 // start at the first company
	stride := 1 // only go one at a time
	start := "0001-01-01" //start 
	end := "3000-01-01" // end, arbitrary date in future means get everything
	// local.RunIssuerCrawl(folder,conn)
	local.RunFormCrawl(folder,start,end,offset,stride,conn)
}
