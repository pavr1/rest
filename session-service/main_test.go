package main

import (
	"net/http"
	"net/http/httptest"
	sharedHealth "shared/health"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// createTestHealthMonitor creates a health monitor for testing with specified state
func createTestHealthMonitor(logger *logrus.Logger, healthy bool) *sharedHealth.HealthMonitor {
	hm := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")
	// Use a mock server or set state directly for testing
	if healthy {
		// Simulate healthy by adding and checking
		hm.SetServiceHealthForTesting("data-service", true)
	}
	return hm
}

func TestHealthCheckEndpoint_Unhealthy(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Health monitor reports data-service as unhealthy (default state)
	healthMonitor := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	healthMonitor.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		healthMonitor:   healthMonitor,
		logger:          logger,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/sessions/p/health", handler.HealthCheck).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/sessions/p/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Data-service unhealthy, expect 503 Service Unavailable
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}
}

func TestHealthCheckEndpoint_Healthy(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Health monitor reports data-service as healthy
	healthMonitor := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	healthMonitor.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")
	healthMonitor.SetServiceHealthForTesting("data-service", true)

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		healthMonitor:   healthMonitor,
		logger:          logger,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/sessions/p/health", handler.HealthCheck).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/sessions/p/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Data-service healthy, expect 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestHealthCheckContentType(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	healthMonitor := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	healthMonitor.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		healthMonitor:   healthMonitor,
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

	healthMonitor := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	healthMonitor.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")
	healthMonitor.SetServiceHealthForTesting("data-service", true)

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		healthMonitor:   healthMonitor,
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

	healthMonitor := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	healthMonitor.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")
	healthMonitor.SetServiceHealthForTesting("data-service", true)

	handler := &MainHTTPHandler{
		sessionsHandler: nil,
		healthMonitor:   healthMonitor,
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
