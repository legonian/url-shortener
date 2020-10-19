package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
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
	// Data type that coming from client
	Url struct {
		Url string `json:"url"`
	}
	// Interface to communicate with main app
	Handler struct {
		DB *sql.DB
	}
)

// Home Page
func (h *Handler) Index(c echo.Context) error {
	return c.File("public/index.html")
	// to test on error:
	// return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
}

// Info Page about URL
func (h *Handler) Info(c echo.Context) error {
	return c.File("public/info.html")
}

// Send new url to database
func (h *Handler) SetRedirectJson(c echo.Context) error {
	var u Url
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, &Data{OK: false})
	}
	urlCode := string(u.Url)
	_, err := url.ParseRequestURI(urlCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Data{OK: false})
	}
	q := fmt.Sprintf("select * from add_url('%s')", urlCode)
	res := getQuery(h.DB, q)
	return c.JSON(http.StatusCreated, res)
}

// Redirect to full URL
func (h *Handler) Redirect(c echo.Context) error {
	short_url := c.Param("short_url")
	cacheData := CheckCache(short_url, h.DB)
	if cacheData != "" {
		return c.Redirect(http.StatusFound, cacheData)
	}
	q := fmt.Sprintf("select * from get_full_url('%s')", short_url)
	res := getQuery(h.DB, q)
	if !res.OK {
		return c.String(http.StatusNotFound, "Shortcut Not Found")
	}
	err := AddToCache(res)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, res.FullURL)
}

// Get info about URL
func (h *Handler) InfoJson(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	q := fmt.Sprintf("select * from get_full_url('%s', 0)", m["url"])
	res := getQuery(h.DB, q)
	return c.JSON(http.StatusOK, res)
}

// Get raw info from database
func getQuery(db *sql.DB, query string) Data {
	rows, err := db.Query(query)
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
