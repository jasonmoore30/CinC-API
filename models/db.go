package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //Need this blank import
)

var db *sql.DB

// InitDB ..
func InitDB(dsn string) {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("sql.Open failed")
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		fmt.Println("Pinging error!")
	}
}
