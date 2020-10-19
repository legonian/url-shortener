package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type (
	// Data type that will be sent to client
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

func (model *DataBaseModel) Init() error {
	sql_url := os.Getenv("DATABASE_URL")
	sqlModel, err := sql.Open("postgres", sql_url)
	model.DB = sqlModel
	if err != nil {
		return err
	}

	if err = model.DB.Ping(); err != nil {
		return err
	}
	log.Println("You connected to database. (DataBaseModel)")
	return nil
}

// Get raw info from database
func (model *DataBaseModel) GetQuery(query string) Data {
	rows, err := model.DB.Query(query)
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
