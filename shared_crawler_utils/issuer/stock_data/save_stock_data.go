package stock_data

import (
	"context"
	"database/sql"
)

// Saves stock data to given db returns true if the issuer was already present
func saveStockData(conn *sql.DB,stockDays []stockDay){
	
	tx,err := conn.BeginTx(context.TODO(),nil)
	if err != nil {
		panic(err)
	}

	stockDaySql := `
	insert into stock_day (ticker,date,close,volume)
	values (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	`
	for _,stockDay := range stockDays {
		_, err := tx.Exec(stockDaySql,stockDay.Ticker,stockDay.Date,stockDay.Close,stockDay.Volume)
		if err != nil {
			panic(err)
		}
	}
                
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

}
