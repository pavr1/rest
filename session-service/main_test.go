package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func TestHealthCheckEndpoint(t *testing.T) {
	// Note: This test will return 503 because data-service is not running
	// In production, data-service must be healthy for session-service to be healthy
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

	// Without data-service running, expect 503 Service Unavailable
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("status code = %d, want %d (data-service not running)", rr.Code, http.StatusServiceUnavailable)
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

	// Response should always be JSON regardless of health status
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
