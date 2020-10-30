package main

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

	"github.com/legonian/url-shortener/database"
)

var (
	// Initial valid URL
	validUrl string = "https://www.example.com/"
	// Valid data from validUrl
	testDataToFill database.Data
	// Array of invalid URLs
	invalidUrlArray []string = []string{"qwerty"}
	// Array of invalid shortcuts codes
	invalidShortcutArray []string = []string{"xxxxxxxxxxxxx"}
)

func TestIndexPage(t *testing.T) {
	r := SetupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET", "/", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
	}
}

func TestStatic(t *testing.T) {
	r := SetupApp()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/public/img/favicon.ico", nil)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
	}

	if body == "" {
		t.Fatalf("no image")
	}
}

func TestCreatingValidLink(t *testing.T) {
	r := SetupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	test_json := fmt.Sprintf(`{"url": "%s"}`, validUrl)
	resp, body := testRequest(t, ts, "POST", "/create", strings.NewReader(test_json))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusCreated)
	}
	if body == "" {
		t.Fatalf("No Info Page")
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatal("Bad MIME type")
	}
	err := json.Unmarshal([]byte(body), &testDataToFill)
	if err != nil {
		t.Fatal(err)
	}

	if !testDataToFill.OK {
		t.Fatal("Wrong json")
	}

	if testDataToFill.FullURL != validUrl {
		t.Fatal("Wrong URL")
	}

	if testDataToFill.ViewsCount != 0 {
		t.Fatal("Views Count not zero after creation")
	}
}

func TestCreatingInvalidLink(t *testing.T) {
	err := database.Init()
	if err != nil {
		log.Fatal(err)
	}
	r := SetupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, invalidUrl := range invalidUrlArray {
		test_json := fmt.Sprintf(`{"url": "%s"}`, invalidUrl)
		resp, _ := testRequest(t, ts, "POST", "/create", strings.NewReader(test_json))
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusBadRequest)
		}
	}
}

func TestInfoPage(t *testing.T) {
	r := SetupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/"+testDataToFill.ShortURL+"/info", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
	}
	if body == "" {
		t.Fatalf("No Info Page")
	}
}

func TestValidLink(t *testing.T) {
	r := SetupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/"+testDataToFill.ShortURL, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
	}
	if body == "" {
		t.Fatalf("No Info Page")
	}
	if resp.Request.URL.String() != testDataToFill.FullURL {
		t.Fatalf("URL is %v, expected %v",
			resp.Request.URL,
			testDataToFill.FullURL,
		)
	}
}

func TestInValidLink(t *testing.T) {
	r := SetupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, invalidShortcut := range invalidShortcutArray {
		resp, _ := testRequest(t, ts, "GET", "/"+invalidShortcut, nil)
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
		}
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	log.SetOutput(ioutil.Discard)

	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
