//Testing main server
package main

import (
	// "net/http"
	// "net/http/httptest"
	"testing"

	// "github.com/stretchr/testify/assert"
)

//TODO

//TestIndexPage
//TestOnCreateURL
//TestOnCreateEmpty
//TestOnUseURL
//TestOnWrongURL

func TestIndexPage(t *testing.T) {
	// router := indexRouter
	// w := httptest.NewRecorder()
	// req, _ := http.NewRequest("GET", "/ping", nil)
	// router.ServeHTTP(w, req)

	// assert.Equal(t, 200, w.Code)
	// assert.Equal(t, "pong", w.Body.String())
	if false {
		t.Error("Failed the test!")
	}
}

func TestOnCreateValidURL(t *testing.T) {}
func TestOnCreateInvalidURL(t *testing.T) {}
func TestOnGoodURL(t *testing.T) {}
func TestOnWrongURL(t *testing.T) {}