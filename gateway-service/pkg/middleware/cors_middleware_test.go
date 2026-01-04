package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware_Headers(t *testing.T) {
	middleware := NewCORSMiddleware(nil)

	handler := middleware.HandleCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Check CORS headers
	tests := []struct {
		header   string
		expected string
	}{
		{"Access-Control-Allow-Origin", "*"},
		{"Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"},
		{"Access-Control-Allow-Credentials", "true"},
		{"Access-Control-Max-Age", "86400"},
	}

	for _, tt := range tests {
		value := w.Header().Get(tt.header)
		if value != tt.expected {
			t.Errorf("%s = %s; want %s", tt.header, value, tt.expected)
		}
	}
}

func TestCORSMiddleware_Preflight(t *testing.T) {
	middleware := NewCORSMiddleware(nil)

	nextCalled := false
	handler := middleware.HandleCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d; want %d", w.Code, http.StatusOK)
	}

	if nextCalled {
		t.Error("next handler should not be called for OPTIONS request")
	}
}

func TestCORSMiddleware_PassThrough(t *testing.T) {
	middleware := NewCORSMiddleware(nil)

	nextCalled := false
	handler := middleware.HandleCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !nextCalled {
		t.Error("next handler should be called for non-OPTIONS request")
	}
}
