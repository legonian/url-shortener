package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestIndexPage(t *testing.T) {
	e, err := initEcho()
	if err != nil {
		t.Fatal(err)
	}
	e.GET("/", handler.Index)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = handler.Index(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("variable value is %v, expected %v", rec.Code, http.StatusOK)
	}
}

func TestInfoPage(t *testing.T) {
	validURL := "https://www.google.com/"
	e, err := initEcho()
	if err != nil {
		t.Fatal(err)
	}

	dataForTest := database.CreateData(validURL)

	e.GET("/:short_url/info", handler.Info)
	e.POST("/:short_url/json", handler.InfoJson)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:short_url/info")
	c.SetParamNames("short_url")
	c.SetParamValues(dataForTest.ShortURL)
	err = handler.Index(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("variable value is %v, expected %v", rec.Code, http.StatusOK)
	}
}

func TestStatic(t *testing.T) {
	e, err := initEcho()
	if err != nil {
		t.Fatal(err)
	}
	e.Static("/public", "public")
	req := httptest.NewRequest(http.MethodGet, "/public/img/favicon.ico", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("variable value is %v, expected %v", rec.Code, http.StatusOK)
	}
	if rec.Body.String() == "" {
		t.Fatalf("no image")
	}
}

func initEcho() (*echo.Echo, error) {
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
	return e, nil
}
