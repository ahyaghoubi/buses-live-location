package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to Database!")
	}

	DB.SetMaxOpenConns(30)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createBusTable := `
		CREATE TABLE IF NOT EXISTS bus (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)
	`

	_, err := DB.Exec(createBusTable)
	if err != nil {
		panic(fmt.Sprintf("Could not create events table: %v", err))
	}

	createAdminTable := `
		CREATE TABLE IF NOT EXISTS admins (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`

	_, err = DB.Exec(createAdminTable)
	if err != nil {
		panic(fmt.Sprintf("Could not create events table: %v", err))
	}

}
