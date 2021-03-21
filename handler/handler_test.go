package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/legonian/url-shortener/database"
)

var (
	// Initial valid URL
	validUrl string = "https://www.example.com/"
	// Valid data from validUrl
	testDataToFill database.Data
)

func TestIndex(t *testing.T) {
	ts, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET", "/", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %v but wanted %v", resp.StatusCode, http.StatusOK)
	}
}

func TestSetRedirectJson(t *testing.T) {
	ts, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	test_json := fmt.Sprintf(`{"url": "%s"}`, validUrl)
	resp, body := testRequest(t, ts, "POST", "/create", strings.NewReader(test_json))
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusCreated)
	}
	if body == "" {
		t.Error("no Info Page")
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("bad MIME type: Content-Type=%v, expect application/json",
			resp.Header.Get("Content-Type"))
	}
	err = json.Unmarshal([]byte(body), &testDataToFill)
	if err != nil {
		t.Error(err)
	}

	if !testDataToFill.OK {
		t.Error("Wrong json")
	}

	if testDataToFill.FullURL != validUrl {
		t.Error("Wrong URL")
	}

	if testDataToFill.ViewsCount != 0 {
		t.Error("Views Count not zero after creation")
	}

	tests := []struct {
		url          string
		expectStatus int
	}{
		{"qwerty", http.StatusBadRequest},
		{"qwerty.", http.StatusBadRequest},
	}

	for _, test := range tests {
		test_json := fmt.Sprintf(`{"url": "%s"}`, test.url)
		resp, _ := testRequest(t, ts, "POST", "/create", strings.NewReader(test_json))
		if resp.StatusCode != test.expectStatus {
			t.Errorf("status code for %s is %d, expected %d",
				test.url, resp.StatusCode, test.expectStatus)
		}
	}
}

func TestInfo(t *testing.T) {
	ts, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/"+testDataToFill.ShortURL+"/info", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
	}
	if body == "" {
		t.Errorf("No Info Page")
	}
}

func TestRedirect(t *testing.T) {
	ts, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/"+testDataToFill.ShortURL, nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
	}
	if body == "" {
		t.Errorf("No Info Page")
	}
	if resp.Request.URL.String() != testDataToFill.FullURL {
		t.Errorf("URL is %v, expected %v",
			resp.Request.URL,
			testDataToFill.FullURL,
		)
	}

	tests := []struct {
		shortCode    string
		expectStatus int
	}{
		{"xxxxxxxxxxxxx", http.StatusNotFound},
	}

	for _, test := range tests {
		resp, _ := testRequest(t, ts, "GET", "/"+test.shortCode, nil)
		if resp.StatusCode != test.expectStatus {
			t.Errorf("status code for %s is %d, expected %d",
				test.shortCode, resp.StatusCode, test.expectStatus)
		}
	}
}

func setupServer() (*httptest.Server, error) {
	if err := database.Init(); err != nil {
		return nil, fmt.Errorf("initializing database error: %s", err)
	}

	if err := SetTemplates("../templates/*"); err != nil {
		return nil, fmt.Errorf("creating temlates error: %s", err)
	}

	r := chi.NewRouter()
	r.Get("/", Index)
	r.Get("/{shortcut}", Redirect)
	r.Get("/{shortcut}/info", Info)

	r.Post("/create", SetRedirectJson)

	ts := httptest.NewServer(r)

	return ts, nil
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	log.SetOutput(ioutil.Discard)

	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Error(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
