package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func TestHealthCheckEndpoint(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		logger:          logger,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/sessions/p/health", handler.HealthCheck).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/sessions/p/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusOK)
	}

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check response structure
	if response["code"] != float64(http.StatusOK) {
		t.Errorf("response code = %v, want %v", response["code"], http.StatusOK)
	}
}

func TestHealthCheckContentType(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		logger:          logger,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/sessions/p/health", handler.HealthCheck).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/sessions/p/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}
}

func TestHealthCheckMethodNotAllowed(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		logger:          logger,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/sessions/p/health", handler.HealthCheck).Methods("GET")

	// Try POST on GET-only endpoint
	req := httptest.NewRequest("POST", "/api/v1/sessions/p/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestSetupRoutes(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		logger:          logger,
	}

	router := mux.NewRouter()
	handler.SetupRoutes(router)

	// Test that routes are registered
	routes := []struct {
		path   string
		method string
	}{
		{"/api/v1/sessions/p/health", "GET"},
		{"/api/v1/sessions/p/login", "POST"},
		{"/api/v1/sessions/p/validate", "POST"},
		{"/api/v1/sessions/logout", "POST"},
	}

	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		match := &mux.RouteMatch{}
		if !router.Match(req, match) {
			t.Errorf("Route %s %s not registered", route.method, route.path)
		}
	}
}
