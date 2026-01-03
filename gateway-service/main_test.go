package main

import (
	"net/http"
	"net/http/httptest"
	sharedHealth "shared/health"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

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

func TestCreateHealthHandler_AllHealthy(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create health monitor with all services healthy
	healthMonitor := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	healthMonitor.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")
	healthMonitor.AddService("session-service", "http://localhost:8087/api/v1/sessions/p/health")
	healthMonitor.SetServiceHealthForTesting("data-service", true)
	healthMonitor.SetServiceHealthForTesting("session-service", true)

	handler := createHealthHandler(healthMonitor, logger)

	req := httptest.NewRequest("GET", "/api/v1/gateway/p/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Should return 200 when all services are healthy
	if w.Code != http.StatusOK {
		t.Errorf("health handler status = %d; want %d", w.Code, http.StatusOK)
	}
}

func TestCreateHealthHandler_ServiceDown(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create health monitor with one service unhealthy
	healthMonitor := sharedHealth.NewHealthMonitor(logger, 10*time.Second)
	healthMonitor.AddService("data-service", "http://localhost:8086/api/v1/data/p/health")
	healthMonitor.AddService("session-service", "http://localhost:8087/api/v1/sessions/p/health")
	healthMonitor.SetServiceHealthForTesting("data-service", true)
	healthMonitor.SetServiceHealthForTesting("session-service", false) // unhealthy

	handler := createHealthHandler(healthMonitor, logger)

	req := httptest.NewRequest("GET", "/api/v1/gateway/p/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Should return 503 when any service is down
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("health handler status = %d; want %d", w.Code, http.StatusServiceUnavailable)
	}
}

func TestCreateProxyHandler_ErrorHandler(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create a proxy handler pointing to a non-existent service
	handler := createProxyHandler("http://localhost:19997", "/api/test", logger)

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
