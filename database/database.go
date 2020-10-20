// Package database provides access to dababase and wrap specific database logic
package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type (
	Data struct {
		OK         bool   `json:"ok"`
		ShortURL   string `json:"short_url"`
		FullURL    string `json:"full_url"`
		ViewsCount int    `json:"views_count"`
	}
	DataBaseModel struct {
		DB *sql.DB
	}
)

var Model DataBaseModel

// Initialize database
func Init() error {
	sql_url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", sql_url)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	Model.DB = db
	log.Println("Database connected.")
	return nil
}

// Get raw info from database
func (this *DataBaseModel) GetQuery(query string) Data {
	rows, err := this.DB.Query(query)
	if err != nil {
		return Data{OK: false}
	}
	defer rows.Close()
	res := make([]Data, 0)
	for rows.Next() {
		url := Data{}
		err := rows.Scan(&url.ShortURL, &url.FullURL, &url.ViewsCount)

		if err != nil {
			res = append(res, Data{OK: false})
		} else {
			url.OK = true
			res = append(res, url)
		}
	}
	if err = rows.Err(); err != nil {
		return Data{OK: false}
	}
	return res[0]
}
