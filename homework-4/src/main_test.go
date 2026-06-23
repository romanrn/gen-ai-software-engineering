package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTimeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "supersecret-hardcoded-key-12345")
	w := httptest.NewRecorder()

	timeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestValidate(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	req.Header.Set("X-Api-Key", "supersecret-hardcoded-key-12345")
	if !validate(req) {
		t.Error("expected valid key to pass")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/time", nil)
	req2.Header.Set("X-Api-Key", "wrong-key")
	if validate(req2) {
		t.Error("expected wrong key to fail")
	}
}