//Testing main server
package main

import (
	"github.com/legonian/url-shortener/database"
	"github.com/legonian/url-shortener/handler"
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
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

func TestCreatingValidURL(t *testing.T) {}
func TestCreatingInvalidURL(t *testing.T) {}
func TestOnGoodURL(t *testing.T) {}
func TestOnWrongURL(t *testing.T) {}

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
