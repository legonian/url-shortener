package main

import (
	"fmt"
	"os"
	"log"
	"net/url"
	"net/http"
	"database/sql"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

// Data type that will be sent to client
type Data struct {
	OK bool `json:"ok"`
	ShortURL string `json:"short_url"`
	FullURL string `json:"full_url"`
	ViewsCount int `json:"views_count"`
}

// Data type that coming from client
type Url struct {
	Url string `json:"url"`
}

// Initialize database
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

// Initialize and run labstack/echo server
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

// Info Page about URL
func infoRouter(c echo.Context) error {
	return c.File("public/info.html")
}

// Send new url to database
func submitRouter(c echo.Context) error {
	var u Url
	//m := echo.Map{}
	if err := c.Bind(&u); err != nil { return err }
	urlCode := string(u.Url)
	_, err := url.ParseRequestURI(urlCode)
	if err != nil {
		return err
	}
	q := fmt.Sprintf("select * from add_url('%s')", urlCode)
	res := get_query(q)
	return c.JSON(http.StatusOK, res)
}

// Redirect to full URL
func redirectByShortURL(c echo.Context) error {
	short_url := c.Param("short_url")
	q := fmt.Sprintf("select * from get_full_url('%s')", short_url)
	res := get_query(q)
	return c.Redirect(http.StatusFound, res.FullURL) // StatusMovedPermanently
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
	q := fmt.Sprintf("select * from get_full_url('%s', 0)", m["url"])
	res := get_query(q)
	return c.JSON(http.StatusOK, res)
}
