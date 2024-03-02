package stock_data

import (
	"database/sql"
	"net/http"

	utils "insiderviz.com/crawlers/shared_crawler_utils"
)
func CrawlStockData(conn *sql.DB,tickers []string, startDate string, num_threads int) {

	thread_guard := make(chan struct{}, num_threads)
	client := &http.Client{}
	for _, ticker := range tickers {
		thread_guard <- struct{}{}
		go func(ticker string) {
			
			utils.ConsumeEODRequest()
			stockData,err := parseStockData(client,ticker,startDate)
			if err != nil {
				panic(err)
			}
			saveStockData(conn,stockData)
			<-thread_guard
		}(ticker)
	}
}