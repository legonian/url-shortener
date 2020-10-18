package database

import (
	"os"
	"log"
	"database/sql"

	_ "github.com/lib/pq"
)

// Initialize database
func Init() (*sql.DB, error) {
	sql_url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", sql_url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("You connected to database.")
	return db, nil
}
