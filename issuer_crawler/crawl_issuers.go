package issuer_crawler

import (
	"database/sql"

	issuer_utils "insiderviz.com/crawlers/utils/issuer"
)

// Crawls all issuers currently in the db, designed to catch changes in companies TODO include some way to delete information that has changed
func CrawlIssuers(conn *sql.DB, num_threads int) {
	tickers := issuer_utils.ReadAllIssuers(conn)
	issuer_utils.CrawlIssuersByCIK(conn, tickers, false, "", num_threads)
}
