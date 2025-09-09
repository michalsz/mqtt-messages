package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	validData := []byte(``)
	req := httptest.NewRequest(http.MethodPost, "/health", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h := HealthCheckHandler{}
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body  %s, got %s", "OK", w.Body)
	}
}
