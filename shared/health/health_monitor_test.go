package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewHealthMonitor(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	hm := NewHealthMonitor(logger, 10*time.Second)

	if hm == nil {
		t.Fatal("NewHealthMonitor() returned nil")
	}

	if hm.interval != 10*time.Second {
		t.Errorf("interval = %v; want %v", hm.interval, 10*time.Second)
	}
}

func TestHealthMonitor_AddService(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	hm := NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("test-service", "http://localhost:8080/health")

	if len(hm.services) != 1 {
		t.Errorf("services count = %d; want 1", len(hm.services))
	}

	svc := hm.services["test-service"]
	if svc == nil {
		t.Fatal("service not added")
	}

	if svc.Name != "test-service" {
		t.Errorf("service name = %s; want test-service", svc.Name)
	}

	if svc.Healthy {
		t.Error("service should start as unhealthy")
	}
}

func TestHealthMonitor_IsServiceHealthy(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	hm := NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("test-service", "http://localhost:8080/health")

	// Initially unhealthy
	if hm.IsServiceHealthy("test-service") {
		t.Error("service should be unhealthy initially")
	}

	// Set to healthy
	hm.SetServiceHealthForTesting("test-service", true)
	if !hm.IsServiceHealthy("test-service") {
		t.Error("service should be healthy after setting")
	}

	// Non-existent service
	if hm.IsServiceHealthy("non-existent") {
		t.Error("non-existent service should return false")
	}
}

func TestHealthMonitor_AreAllServicesHealthy(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	hm := NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("service-1", "http://localhost:8080/health")
	hm.AddService("service-2", "http://localhost:8081/health")

	// Initially all unhealthy
	if hm.AreAllServicesHealthy() {
		t.Error("should be false when services are unhealthy")
	}

	// One healthy
	hm.SetServiceHealthForTesting("service-1", true)
	if hm.AreAllServicesHealthy() {
		t.Error("should be false when one service is unhealthy")
	}

	// Both healthy
	hm.SetServiceHealthForTesting("service-2", true)
	if !hm.AreAllServicesHealthy() {
		t.Error("should be true when all services are healthy")
	}
}

func TestHealthMonitor_GetServiceStatus(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	hm := NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("test-service", "http://localhost:8080/health")
	hm.SetServiceHealthForTesting("test-service", true)

	status := hm.GetServiceStatus("test-service")
	if status == nil {
		t.Fatal("GetServiceStatus returned nil")
	}

	if status.Name != "test-service" {
		t.Errorf("name = %s; want test-service", status.Name)
	}

	if !status.Healthy {
		t.Error("status should be healthy")
	}

	// Non-existent service
	status = hm.GetServiceStatus("non-existent")
	if status != nil {
		t.Error("non-existent service should return nil")
	}
}

func TestHealthMonitor_GetAllServicesStatus(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	hm := NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("service-1", "http://localhost:8080/health")
	hm.AddService("service-2", "http://localhost:8081/health")

	all := hm.GetAllServicesStatus()
	if len(all) != 2 {
		t.Errorf("services count = %d; want 2", len(all))
	}
}

func TestHealthMonitor_CheckService_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create a test server that returns 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hm := NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("test-service", server.URL)

	// Manually check the service
	svc := hm.services["test-service"]
	hm.checkService(svc)

	if !hm.IsServiceHealthy("test-service") {
		t.Error("service should be healthy after successful check")
	}
}

func TestHealthMonitor_CheckService_Failure(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create a test server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	hm := NewHealthMonitor(logger, 10*time.Second)
	hm.AddService("test-service", server.URL)

	// Manually check the service
	svc := hm.services["test-service"]
	hm.checkService(svc)

	if hm.IsServiceHealthy("test-service") {
		t.Error("service should be unhealthy after failed check")
	}
}

