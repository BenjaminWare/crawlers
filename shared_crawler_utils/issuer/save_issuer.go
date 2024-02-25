package issuer

import (
	"context"
	"database/sql"

	. "insiderviz.com/crawlers/shared_crawler_utils"
)

// Saves issuer to given db returns true if the issuer was already present
func saveIssuer(conn *sql.DB,issuer Issuer) bool{
	
	tx,err := conn.BeginTx(context.TODO(),nil)
	if err != nil {
		panic(err)
	}

	issuerSql := `
	replace into issuer (cik, name, sic, sic_description, ein, state_of_incorporation,fiscal_year_end, phone, sector, industry)
	values (?, ?, ?, ?, ?, ?, ?, ?, ?,?)
	`

	result, err := tx.Exec(issuerSql,issuer.Cik,issuer.Name,issuer.Sic,issuer.Sic_description,issuer.Ein,issuer.State_of_incorporation,issuer.Fiscal_year_end,issuer.Phone,issuer.Sector,issuer.Industry)
	if err != nil {
		panic(err)
	}

	rows,err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	// True when the issuer was already present thus rows affected is 2 one deleted one added
	return rows == 2
}

func saveTickers(conn *sql.DB,issuer Issuer) {

	tx,err := conn.BeginTx(context.TODO(),nil)
	if err != nil {
		panic(err)
	}
	
	tickersSql := `
	replace into ticker (cik,ticker)
	values (?,?)
	`
	for _,ticker := range issuer.Tickers {
		_, err = tx.Exec(tickersSql,issuer.Cik,ticker)
			if err != nil {
				panic(err)
			}
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}