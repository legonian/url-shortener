package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"database/sql"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

type Data struct {
	OK bool `json:"ok"`
	ShortURL string `json:"short_url"`
	FullURL string `json:"full_url"`
	ViewsCount int `json:"views_count"`
}

func init() {
	var err error
	sql_url := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", sql_url)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	e := echo.New()

  e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/public", "public")

	e.GET("/", indexRouter)
	e.GET("/:short_url", redirectByShortURL)
	e.GET("/:short_url/info", infoRouter)

	e.POST("/submit", submitRouter)
	e.POST("/api/:short_url", apiRouter)
	
	e.Logger.Fatal(e.Start(":" + port))
}

// Routes

// Home Page
func indexRouter(c echo.Context) error {
	return c.File("public/index.html")
}

// Info about short URL
func infoRouter(c echo.Context) error {
	return c.File("public/info.html")
}

func submitRouter(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	q := fmt.Sprintf("select * from add_url('%s')", m["url"])
	log.Print(q)
	res := get_query(q)
	return c.JSON(http.StatusOK, res)
}

// Redirect to full URL
func redirectByShortURL(c echo.Context) error {
	short_url := c.Param("short_url")
	q := fmt.Sprintf("select * from get_full_url('%s')", short_url)
	res := get_query(q)
	log.Print(res.FullURL)
	return c.Redirect(http.StatusMovedPermanently, res.FullURL) // http.StatusMovedPermanently
}

// Functions

// Get raw info from database
func get_query(query string) Data{
	rows, err := db.Query(query)
	if err != nil {
		return Data{OK:false}
	}
	defer rows.Close()
	res := make([]Data, 0)
	for rows.Next() {
		url := Data{}
		err := rows.Scan(&url.ShortURL, &url.FullURL, &url.ViewsCount)
		
		if err != nil {
			res = append(res, Data{OK:false})
		} else {
			url.OK = true
			res = append(res, url)
		}
		
	}
	if err = rows.Err(); err != nil {
		return Data{OK:false}
	}
	return res[0]
}

// Get info about URL
func apiRouter(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	log.Print(m)
	q := fmt.Sprintf("select * from get_full_url('%s')", m["url"])
	res := get_query(q)
	log.Print(res)
	return c.JSON(http.StatusOK, res)
}
