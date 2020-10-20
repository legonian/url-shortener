//Testing main server
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/legonian/url-shortener/database"
	"github.com/legonian/url-shortener/handler"
)

type (
	Data struct {
		OK         bool   `json:"ok"`
		ShortURL   string `json:"short_url"`
		FullURL    string `json:"full_url"`
		ViewsCount int    `json:"views_count"`
	}
)

var (
	validURL   string = "https://www.google.com/"
	invalidURL string = "qwerty"
	d          Data
)

func TestIndexPage(t *testing.T) {
	e, err := init_app()
	expect(t, err, nil)
	e.GET("/", handler.Index)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, rec)
	err = handler.Index(c)
	expect(t, err, nil)
}

func TestCreatingValidURL(t *testing.T) {
	e, err := init_app()
	expect(t, err, nil)
	e.POST("/create", handler.SetRedirectJson)

	test_json := fmt.Sprintf(`{"url": "%s"}`, validURL)
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = handler.SetRedirectJson(c)
	expect(t, err, nil)
	expect(t, rec.Code, http.StatusCreated)
	err = json.NewDecoder(rec.Body).Decode(&d)
	expect(t, err, nil)
}

func TestCreatingInvalidURL(t *testing.T) {
	e, err := init_app()
	expect(t, err, nil)
	e.POST("/create", handler.SetRedirectJson)

	test_json := fmt.Sprintf(`{"url": "%s"}`, invalidURL)
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = handler.SetRedirectJson(c)
	expect(t, err, nil)
	expect(t, rec.Code, http.StatusBadRequest)
}

func TestOnGoodURL(t *testing.T) {
	e, err := init_app()
	expect(t, err, nil)
	e.GET("/:short_url", handler.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:short_url")
	c.SetParamNames("short_url")
	c.SetParamValues(d.ShortURL)
	err = handler.Redirect(c)
	redirectedURL := rec.HeaderMap["Location"][0]

	expect(t, err, nil)
	expect(t, rec.Code, http.StatusFound)
	expect(t, redirectedURL, validURL)
}

func TestOnWrongURL(t *testing.T) {
	e, err := init_app()
	expect(t, err, nil)
	e.GET("/:short_url", handler.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/xxxxxxxxxxxxx", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = handler.Redirect(c)

	expect(t, err, nil)
	expect(t, rec.Code, http.StatusNotFound)
}

func init_app() (*echo.Echo, error) {
	log.SetOutput(ioutil.Discard)

	err := database.Init()
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "path:${uri} | ${method} method to ${status} | t=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Static("/public", "public")

	return e, err
}

func expect(t *testing.T, varToTest interface{}, expected interface{}) {
	if varToTest != expected {
		t.Fatalf("variable value is %v, expected %v", varToTest, expected)
	}
}
