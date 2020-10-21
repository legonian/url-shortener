// Package database provides access to dababase and wrap any specific database
// logic
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

var (
	Model DataBaseModel
)

const (
	IsViewed  int = 1
	NotViewed int = 0
)

// Initialize database by connecting through DATABASE_URL env variable and
// specify database/sql driver
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

// Add original URL to database and get data with type of Data that has shortcut
// code and can be used to redirect to original URL
func CreateData(fullUrlToSave string) Data {
	stmt, err := Model.DB.Prepare("select * from add_url($1)")
	if err != nil {
		log.Println(err)
		return Data{OK: false}
	}
	defer stmt.Close()

	data := Data{OK: true}
	row := stmt.QueryRow(fullUrlToSave)
	err = row.Scan(&data.ShortURL, &data.FullURL, &data.ViewsCount)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return Data{OK: false}
	}
	return data
}

// Get data from database with Data type by sending shortUrl code that represent
// data in database and viewIncrease that used to increament views counter if
// needed to show that data was viewed
func GetData(shortUrl string, viewIncrease int) Data {
	stmt, err := Model.DB.Prepare("select * from get_full_url($1,$2)")
	if err != nil {
		return Data{OK: false}
	}
	defer stmt.Close()

	data := Data{OK: true}
	row := stmt.QueryRow(shortUrl, viewIncrease)
	err = row.Scan(&data.ShortURL, &data.FullURL, &data.ViewsCount)
	if err != nil && err != sql.ErrNoRows {
		return Data{OK: false}
	}
	return data
}
