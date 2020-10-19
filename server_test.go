//Testing main server
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

var (
	validURL   string = "https://www.google.com/"
	invalidURL string = "qwerty"
	d          Data
)

func TestIndexPage(t *testing.T) {
	db, e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	h := &handler.Handler{DB: db}
	e.GET("/", h.Index)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, rec)
	err = h.Index(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreatingValidURL(t *testing.T) {
	db, e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	h := &handler.Handler{DB: db}
	e.POST("/create", h.SetRedirectJson)

	test_json := fmt.Sprintf(`{"url": "%s"}`, validURL)
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = h.SetRedirectJson(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatal("Not succeeded")
	}
	err = json.NewDecoder(rec.Body).Decode(&d)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreatingInvalidURL(t *testing.T) {
	db, e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	h := &handler.Handler{DB: db}
	e.POST("/create", h.SetRedirectJson)

	test_json := fmt.Sprintf(`{"url": "%s"}`, invalidURL)
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = h.SetRedirectJson(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatal("Probably succeeded")
	}
}

func TestOnGoodURL(t *testing.T) {
	db, e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	h := &handler.Handler{DB: db}
	e.GET("/:short_url", h.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:short_url")
	c.SetParamNames("short_url")
	c.SetParamValues(d.ShortURL)

	err = h.Redirect(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusFound {
		t.Fatal(rec.Body)
	}
	redirectedURL := rec.HeaderMap["Location"][0]
	if redirectedURL != validURL {
		t.Fatal(rec.HeaderMap)
	}
}

func TestOnWrongURL(t *testing.T) {
	db, e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	h := &handler.Handler{DB: db}
	e.GET("/:short_url", h.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/xxxxxxxxxxxxx", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = h.Redirect(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatal(rec.Code)
	}
}

func init_app() (*sql.DB, *echo.Echo, error) {
	db, err := database.Init()
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
