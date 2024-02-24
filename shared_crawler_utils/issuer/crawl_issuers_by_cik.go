package issuer

import (
	"database/sql"
)

func CrawlIssuersByCIK(conn *sql.DB, ciks []string) {
	for _, cik := range ciks {
		issuer := crawlIssuerJSON(cik)
		// Uses 0 padded version of cik
		issuer.Cik = cik
		saveIssuer(conn, issuer)
	}
}