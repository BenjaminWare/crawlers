package main

import (
	"database/sql"

	. "insiderviz.com/crawlers/live_crawler"
	. "insiderviz.com/crawlers/local_crawler"
	. "insiderviz.com/crawlers/shared_crawler_functions"
)

// This file can run any crawler, should mainly be used for local testing as crawlers should be deployed indiviually
func main() {
	// Creates connection to local db called insiderviz-crawler2 with username and password "root"
	conn := CreateMySQLConnection("root:root@tcp(127.0.0.1:3306)/insiderviz-crawler2") 
	defer conn.Close()

	test_live(conn)
	// test_local(conn)
}

func test_live(conn *sql.DB) {
	RunLiveCrawler(conn)
}

func test_local(conn *sql.DB) {
	folder := "./submissions/2024-01-27"

	offset := 0 // start at the first company
	stride := 1 // only go one at a time
	start := "2023-01-16" //start 
	end := "3000-01-01" // end, arbitrary date in future means get everything

	RunLocalCrawler(folder,start,end,offset,stride,conn)
}