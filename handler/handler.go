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
	// Data type that coming from client to create Data
	Url struct {
		Url string `json:"url"`
	}
)

// Open Index html page
func Index(c echo.Context) error {
	return c.File("public/index.html")
	// to test on error:
	// return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
}

// Open Info html page
func Info(c echo.Context) error {
	return c.File("public/info.html")
}

// Parse POST request, send given URL to database, after getting data with
// shortcut code send it to client in json
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
	res := database.CreateData(urlCode)
	return c.JSON(http.StatusCreated, res)
}

// Parse GET request, and after getting data about given given shortcut
// link from database(or cache) redirect client to full link
func Redirect(c echo.Context) error {
	urlCode := c.Param("short_url")
	cacheData := database.CheckCache(urlCode, database.IsViewed)
	if cacheData.OK && cacheData.FullURL != "" {
		return c.Redirect(http.StatusFound, cacheData.FullURL)
	}
	res := database.GetData(urlCode, database.IsViewed)
	if !res.OK {
		return c.String(http.StatusNotFound, "Shortcut Not Found")
	}
	database.AddCache(res)
	return c.Redirect(http.StatusFound, res.FullURL)
}

// Parse POST request, and after getting data about given given shortcut
// link from database(or cache) send it to client in json
func InfoJson(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	urlCode := fmt.Sprintf("%s", m["url"])
	cacheData := database.CheckCache(urlCode, database.NotViewed)
	if cacheData.OK {
		return c.JSON(http.StatusOK, cacheData)
	}
	res := database.GetData(urlCode, database.NotViewed)
	if !res.OK {
		return c.String(http.StatusNotFound, "Shortcut Not Found")
	}
	database.AddCache(res)
	return c.JSON(http.StatusOK, res)
}
