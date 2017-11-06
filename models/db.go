package models

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

// InitDB ..
func InitDB(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("sql.Open failed")
	}

	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
