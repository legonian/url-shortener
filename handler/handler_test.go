package handler

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
)

var (
	dataToFillThroughTest         Data
	validLinkToTest               string   = "https://www.google.com/"
	invalidArrayOfLinkToTest      []string = []string{"qwerty"} // Full path Links
	invalidArrayOfShortLinkToTest []string = []string{"qwerty"} // Shortcut Links
)

func TestCreatingValidLink(t *testing.T) {
	e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	e.POST("/create", SetRedirectJson)

	test_json := fmt.Sprintf(`{"url": "%s"}`, validLinkToTest)
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = SetRedirectJson(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("variable value is %v, expected %v", rec.Code, http.StatusCreated)
	}
	err = json.NewDecoder(rec.Body).Decode(&dataToFillThroughTest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreatingInvalidLink(t *testing.T) {
	e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	e.POST("/create", SetRedirectJson)

	for _, invalidLink := range invalidArrayOfLinkToTest {
		test_json := fmt.Sprintf(`{"url": "%s"}`, invalidLink)
		req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test_json))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = SetRedirectJson(c)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("variable value is %v, expected %v", rec.Code, http.StatusBadRequest)
		}
	}

}

func TestValidLink(t *testing.T) {
	e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	e.GET("/:short_url", Redirect)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:short_url")
	c.SetParamNames("short_url")
	c.SetParamValues(dataToFillThroughTest.ShortURL)
	err = Redirect(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusFound {
		t.Fatalf("variable value is %v, expected %v",
			rec.Code,
			http.StatusFound)
	}
	if rec.HeaderMap["Location"][0] != validLinkToTest {
		t.Fatalf("variable value is %v, expected %v",
			rec.HeaderMap["Location"][0],
			validLinkToTest)
	}
}

func TestInValidLink(t *testing.T) {
	e, err := init_app()
	if err != nil {
		t.Fatal(err)
	}
	e.GET("/:short_url", Redirect)

	for _, invalidShortLink := range invalidArrayOfShortLinkToTest {
		req := httptest.NewRequest(http.MethodGet, "/"+invalidShortLink, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Redirect(c)

		if err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusNotFound {
			t.Fatalf("variable value is %v, expected %v", rec.Code, http.StatusNotFound)
		}
	}
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
