package stock_data

import (
	"database/sql"
	"net/http"
	"sync"

	utils "insiderviz.com/crawlers/utils"
)
func CrawlStockData(conn *sql.DB,tickers []string, startDate string, num_threads int) {

	thread_guard := make(chan struct{}, num_threads)
	client := &http.Client{}
	var wg sync.WaitGroup
	for _, ticker := range tickers {
		thread_guard <- struct{}{}
		wg.Add(1)
		go func(ticker string) {
			defer wg.Done()
			utils.ConsumeEODRequest()
			stockData,err := parseStockData(client,ticker,startDate)

			if err != nil {
				panic(err)
			}
			saveStockData(conn,stockData)
			<-thread_guard
			println("Finished:",ticker)
		}(ticker)
	}
	wg.Wait()
}