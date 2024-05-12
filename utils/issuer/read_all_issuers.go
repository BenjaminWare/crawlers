package issuer

import "database/sql"

// Helper function to get all issuers currently in the mysql db @conn and return them as a slice
func ReadAllIssuers(conn *sql.DB) []string {
	tickers := make([]string, 0)
	tickerSQL := `select ticker from ticker`
	result, err := conn.Query(tickerSQL)
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var ticker string
		result.Scan(&ticker)
		tickers = append(tickers, ticker)
	}
	return tickers
}
