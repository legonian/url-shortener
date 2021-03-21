package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexPage(t *testing.T) {
	r := setupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET", "/", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("home page status code is %d, expected %d",
			resp.StatusCode, http.StatusOK)
	}
}

func TestStatic(t *testing.T) {
	r := setupApp()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/public/img/favicon.ico", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode is %v, expected %v", resp.StatusCode, http.StatusOK)
	}
	if body == "" {
		t.Errorf("no image")
	}
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
