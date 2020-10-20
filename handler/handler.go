// Package handler provides functions that represent routing logic
package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/legonian/url-shortener/database"
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
)

// Home Page
func Index(c echo.Context) error {
	return c.File("public/index.html")
	// to test on error:
	// return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
}

// Info Page about URL
func Info(c echo.Context) error {
	return c.File("public/info.html")
}

// Send new url to database
func SetRedirectJson(c echo.Context) error {
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
	res := database.Model.GetQuery(q)
	return c.JSON(http.StatusCreated, res)
}

// Redirect to full URL
func Redirect(c echo.Context) error {
	urlCode := c.Param("short_url")
	cacheData := database.CheckCache(urlCode, true)
	if cacheData.OK && cacheData.FullURL != "" {
		return c.Redirect(http.StatusFound, cacheData.FullURL)
	}

	q := fmt.Sprintf("select * from get_full_url('%s')", urlCode)
	res := database.Model.GetQuery(q)
	if !res.OK {
		return c.String(http.StatusNotFound, "Shortcut Not Found")
	}
	err := database.AddCache(res)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, res.FullURL)
}

// Get info about URL
func InfoJson(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	urlCode := fmt.Sprintf("%s", m["url"])
	cacheData := database.CheckCache(urlCode, false)
	if cacheData.OK {
		return c.JSON(http.StatusOK, cacheData)
	}

	q := fmt.Sprintf("select * from get_full_url('%s', 0)", urlCode)
	res := database.Model.GetQuery(q)
	err := database.AddCache(res)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
