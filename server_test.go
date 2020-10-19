//Testing main server
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/legonian/url-shortener/database"
	"github.com/legonian/url-shortener/handler"
	_ "github.com/lib/pq"
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

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestIndexPage(t *testing.T) {
	db, e, err := init_app()
	expect(t, err, nil)
	h := &handler.Handler{DB: db}
	e.GET("/", h.Index)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, rec)
	err = h.Index(c)
	expect(t, err, nil)
}

func TestCreatingValidURL(t *testing.T) {
	db, e, err := init_app()
	expect(t, err, nil)
	h := &handler.Handler{DB: db}
	e.POST("/create", h.SetRedirectJson)

	test_json := fmt.Sprintf(`{"url": "%s"}`, validURL)
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = h.SetRedirectJson(c)
	expect(t, err, nil)
	expect(t, rec.Code, http.StatusCreated)
	err = json.NewDecoder(rec.Body).Decode(&d)
	expect(t, err, nil)
}

func TestCreatingInvalidURL(t *testing.T) {
	db, e, err := init_app()
	expect(t, err, nil)
	h := &handler.Handler{DB: db}
	e.POST("/create", h.SetRedirectJson)

	test_json := fmt.Sprintf(`{"url": "%s"}`, invalidURL)
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = h.SetRedirectJson(c)
	expect(t, err, nil)
	expect(t, rec.Code, http.StatusBadRequest)
}

func TestOnGoodURL(t *testing.T) {
	db, e, err := init_app()
	expect(t, err, nil)
	h := &handler.Handler{DB: db}
	e.GET("/:short_url", h.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:short_url")
	c.SetParamNames("short_url")
	c.SetParamValues(d.ShortURL)
	err = h.Redirect(c)
	redirectedURL := rec.HeaderMap["Location"][0]

	expect(t, err, nil)
	expect(t, rec.Code, http.StatusFound)
	expect(t, redirectedURL, validURL)
}

func TestOnWrongURL(t *testing.T) {
	db, e, err := init_app()
	expect(t, err, nil)
	h := &handler.Handler{DB: db}
	e.GET("/:short_url", h.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/xxxxxxxxxxxxx", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = h.Redirect(c)

	expect(t, err, nil)
	expect(t, rec.Code, http.StatusNotFound)
}

func init_app() (*database.DataBaseModel, *echo.Echo, error) {
	db := &database.Model
	err := db.Init()
	if err != nil {
		return nil, nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		err = errors.New("$PORT must be set")
		return nil, nil, err
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/public", "public")

	return db, e, err
}

func expect(t *testing.T, varToTest interface{}, expected interface{}) {
	if varToTest != expected {
		t.Fatalf("variable value is %v, expected %v", varToTest, expected)
	}
}
