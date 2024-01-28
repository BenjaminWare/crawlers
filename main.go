package main

import (
	"database/sql"

	live "insiderviz.com/live_crawler"
	local "insiderviz.com/local_crawler"
	crawler_utils "insiderviz.com/shared_crawler_functions"
)

// This file can run any crawler, should mainly be used for local testing as crawlers should be deployed indiviually
func main() {
	// Creates connection to local db called insiderviz-crawler2 with username and password "root"
	conn := crawler_utils.CreateMySQLConnection("root:root@tcp(127.0.0.1:3306)/insiderviz-crawler2") 
	defer conn.Close()

	test_live(conn)
	// test_local(conn)
}

func test_live(conn *sql.DB) {
	live.RunLiveCrawler(conn)
}

func test_local(conn *sql.DB) {
	folder := "./submissions/2024-01-27"

	offset := 0 // start at the first company
	stride := 1 // only go one at a time
	start := "2023-01-16" //start 
	end := "3000-01-01" // end, arbitrary date in future means get everything

	local.RunLocalCrawler(folder,start,end,offset,stride,conn)
}