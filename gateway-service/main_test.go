package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBoolToHealth(t *testing.T) {
	tests := []struct {
		input    bool
		expected string
	}{
		{true, "healthy"},
		{false, "unhealthy"},
	}

	for _, tt := range tests {
		result := boolToHealth(tt.input)
		if result != tt.expected {
			t.Errorf("boolToHealth(%v) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	if id1 == "" {
		t.Error("generateRequestID() returned empty string")
	}

	if len(id1) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("generateRequestID() length = %d; want 32", len(id1))
	}

	if id1 == id2 {
		t.Error("generateRequestID() should return unique IDs")
	}
}

func TestCheckServiceHealth_InvalidURL(t *testing.T) {
	result := checkServiceHealth("http://invalid-url-that-does-not-exist:9999/health", nil)
	if result {
		t.Error("checkServiceHealth() should return false for invalid URL")
	}
}

func TestCreateHealthHandler_ServiceDown(t *testing.T) {
	// Create handler with non-existent services
	handler := createHealthHandler(
		"http://localhost:19999", // non-existent session service
		"http://localhost:19998", // non-existent data service
		nil,
	)

	req := httptest.NewRequest("GET", "/api/v1/gateway/p/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Should return 503 when services are down
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("health handler status = %d; want %d", w.Code, http.StatusServiceUnavailable)
	}
}

func TestCreateProxyHandler_ErrorHandler(t *testing.T) {
	// Create a proxy handler pointing to a non-existent service
	handler := createProxyHandler("http://localhost:19997", "/api/test", nil)

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Should return 502 Bad Gateway
	if w.Code != http.StatusBadGateway {
		t.Errorf("proxy handler status = %d; want %d", w.Code, http.StatusBadGateway)
	}

	// Should have Content-Type application/json
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s; want application/json", contentType)
	}
}
