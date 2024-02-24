package issuer

import (
	"context"
	"database/sql"

	. "insiderviz.com/crawlers/shared_crawler_utils"
)


func saveIssuer(conn *sql.DB,issuer Issuer) {
	
	tx,err := conn.BeginTx(context.TODO(),nil)
	if err != nil {
		panic(err)
	}

	issuerSql := `
	insert into issuer (cik, name, sic, sic_description, ein, state_of_incorporation,fiscal_year_end, phone, sector, industry)
	values (?, ?, ?, ?, ?, ?, ?, ?, ?,?)
	ON DUPLICATE KEY UPDATE cik=cik
	`

	_, err = tx.Exec(issuerSql,issuer.Cik,issuer.Name,issuer.Sic,issuer.Sic_description,issuer.Ein,issuer.State_of_incorporation,issuer.Fiscal_year_end,issuer.Phone,issuer.Sector,issuer.Industry)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

