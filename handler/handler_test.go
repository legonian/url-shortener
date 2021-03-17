package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	if err := SetTemplates("../templates/*"); err != nil {
		t.Fatalf("creating temlates error: %s", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	Index(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d but wanted %d", rec.Code, http.StatusOK)
	}
}
