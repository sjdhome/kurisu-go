package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDatabase() *sql.DB {
	return db
}

func Init() error {
	_db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		return fmt.Errorf("Opening database: %w", err)
	}
	db = _db
	return nil
}
