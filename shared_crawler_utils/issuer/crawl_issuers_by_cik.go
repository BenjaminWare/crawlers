package issuer

import (
	"database/sql"
	"fmt"
	"sync"

	stock_data_utils "insiderviz.com/crawlers/shared_crawler_utils/issuer/stock_data"
)

/*
	Crawls the given ciks into the given conn population the issuer and tickers tables
	Can crawl stock data if crawl_stock_data = True getting all days after stock_data_start_date, stock_data_start_date does nothing when crawl_stock_data is false
	When crawling stock_data extra threads will be created to get stock_data so consider reducing num_threads
*/
func CrawlIssuersByCIK(conn *sql.DB, ciks []string,crawl_stock_data bool,stock_data_start_date string,num_threads int) {
	// Controls number of threads spun up
	thread_guard := make(chan struct{} , num_threads)
	// Makes sure the function waits until all routines finish
	var wg sync.WaitGroup

	for i, cik := range ciks {
		thread_guard <- struct{}{}
		wg.Add(1)
		
		go func(cik string, count int) {
			defer wg.Done()
			issuer := parseIssuerJSON(cik)

			// Uses 0 padded version of cik
			issuer.Cik = cik
			saveIssuer(conn, issuer)
			saveTickers(conn,issuer)
			if crawl_stock_data {
				stock_data_utils.CrawlStockData(conn,issuer.Tickers,stock_data_start_date,1)
			}
	
			fmt.Printf("Finished%d/%d: %s\t%s\n",count+1,len(ciks),issuer.Name,issuer.Cik)
			<-thread_guard
		}(cik,i)


	}

	wg.Wait()
}