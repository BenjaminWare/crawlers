package stock_day_crawler

import (
	"database/sql"
	"fmt"
	"time"

	issuer_utils "insiderviz.com/crawlers/utils/issuer"
	stock_data_utils "insiderviz.com/crawlers/utils/issuer/stock_data"
)

func CrawlTodaysStockDay(conn *sql.DB, num_threads int) {
	year, month, day := time.Now().Date()
	today := fmt.Sprintf("%d-%d-%d", year, int(month), day)
	tickers := issuer_utils.ReadAllIssuers(conn)
	stock_data_utils.CrawlStockData(conn, tickers, today, num_threads)
}
