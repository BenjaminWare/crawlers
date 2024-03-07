package utils

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Returns a reference to the db
func CreateMySQLConnection(DSN string) *sql.DB {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}

	return db
}
