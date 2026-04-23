package dbmanager

import (
	"database/sql"
	"log"
)

func NewGestotorDb(dataBasePath string) *sql.DB {
	db, err := sql.Open("sqlite3", dataBasePath) // "./database/app.db"
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
