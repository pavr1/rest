package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"data-service/pkg/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// MockDatabaseHandler is a mock implementation of DatabaseHandler for testing
type MockDatabaseHandler struct {
	HealthCheckFunc func() error
	GetStatsFunc    func() sql.DBStats
	QueryFunc       func(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowFunc    func(query string, args ...interface{}) *sql.Row
	ExecFunc        func(query string, args ...interface{}) (sql.Result, error)
}

func (m *MockDatabaseHandler) Connect() error                               { return nil }
func (m *MockDatabaseHandler) Close() error                                 { return nil }
func (m *MockDatabaseHandler) Ping() error                                  { return nil }
func (m *MockDatabaseHandler) BeginTx(ctx context.Context) (*sql.Tx, error) { return nil, nil }
func (m *MockDatabaseHandler) CommitTx(tx *sql.Tx) error                    { return nil }
func (m *MockDatabaseHandler) RollbackTx(tx *sql.Tx) error                  { return nil }
func (m *MockDatabaseHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}
func (m *MockDatabaseHandler) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}
func (m *MockDatabaseHandler) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (m *MockDatabaseHandler) Prepare(query string) (*sql.Stmt, error) { return nil, nil }
func (m *MockDatabaseHandler) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}
func (m *MockDatabaseHandler) GetDB() *sql.DB    { return nil }
func (m *MockDatabaseHandler) IsConnected() bool { return true }

func (m *MockDatabaseHandler) HealthCheck() error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc()
	}
	return nil
}

func (m *MockDatabaseHandler) GetStats() sql.DBStats {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc()
	}
	return sql.DBStats{
		OpenConnections: 5,
		InUse:           2,
		Idle:            3,
	}
}

func (m *MockDatabaseHandler) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, nil
}

func (m *MockDatabaseHandler) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(query, args...)
	}
	return nil
}

func (m *MockDatabaseHandler) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}

// Helper function to create test handler
func createTestHandler(mockDB *MockDatabaseHandler) *Handler {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Suppress logs during tests

	config := &database.Config{
		Host:   "localhost",
		Port:   5432,
		DBName: "test_db",
	}

	return &Handler{
		db:     mockDB,
		config: config,
		logger: logger,
	}
}

// TestNewHandler tests handler creation
func TestNewHandler(t *testing.T) {
	mockDB := &MockDatabaseHandler{}
	logger := logrus.New()
	config := &database.Config{
		Host:   "localhost",
		Port:   5432,
		DBName: "test_db",
	}

	handler := NewHandler(mockDB, config, logger)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.db != mockDB {
		t.Error("Expected db to be set correctly")
	}

	if handler.config != config {
		t.Error("Expected config to be set correctly")
	}

	if handler.logger != logger {
		t.Error("Expected logger to be set correctly")
	}
}

// TestSetupRoutes tests that routes are configured correctly
func TestSetupRoutes(t *testing.T) {
	mockDB := &MockDatabaseHandler{}
	handler := createTestHandler(mockDB)

	router := mux.NewRouter()
	handler.SetupRoutes(router)

	// Test that routes exist by making requests
	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"GET", "/api/v1/data/p/health"},
		{"GET", "/api/v1/data/p/stats"},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		match := &mux.RouteMatch{}
		if !router.Match(req, match) {
			t.Errorf("Expected route %s %s to be registered", tc.method, tc.path)
		}
	}
}

// TestRootHandler tests the root endpoint
func TestRootHandler(t *testing.T) {
	mockDB := &MockDatabaseHandler{}
	handler := createTestHandler(mockDB)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.RootHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] != "Bar-Restaurant Data Service is running" {
		t.Errorf("Unexpected message: %s", response["message"])
	}
}

// TestHealthCheck_Healthy tests health check when database is healthy
func TestHealthCheck_Healthy(t *testing.T) {
	mockDB := &MockDatabaseHandler{
		HealthCheckFunc: func() error {
			return nil
		},
		GetStatsFunc: func() sql.DBStats {
			return sql.DBStats{
				OpenConnections: 10,
				InUse:           3,
				Idle:            7,
			}
		},
	}
	handler := createTestHandler(mockDB)

	req := httptest.NewRequest("GET", "/api/v1/data/p/health", nil)
	rr := httptest.NewRecorder()

	handler.HealthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "data-service" {
		t.Errorf("Expected service 'data-service', got '%s'", response["service"])
	}

	if response["message"] != "Database ping check passed" {
		t.Errorf("Unexpected message: %s", response["message"])
	}
}

// TestHealthCheck_Unhealthy tests health check when database is unhealthy
func TestHealthCheck_Unhealthy(t *testing.T) {
	mockDB := &MockDatabaseHandler{
		HealthCheckFunc: func() error {
			return errors.New("database connection failed")
		},
	}
	handler := createTestHandler(mockDB)

	req := httptest.NewRequest("GET", "/api/v1/data/p/health", nil)
	rr := httptest.NewRecorder()

	handler.HealthCheck(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, rr.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "unhealthy" {
		t.Errorf("Expected status 'unhealthy', got '%s'", response["status"])
	}

	if response["message"] != "Database ping check failed" {
		t.Errorf("Unexpected message: %s", response["message"])
	}

	if response["error"] != "database connection failed" {
		t.Errorf("Expected error message, got '%s'", response["error"])
	}
}

// TestStatsEndpoint tests the stats endpoint
func TestStatsEndpoint(t *testing.T) {
	mockDB := &MockDatabaseHandler{
		GetStatsFunc: func() sql.DBStats {
			return sql.DBStats{
				OpenConnections: 15,
				InUse:           5,
				Idle:            10,
				WaitCount:       2,
			}
		},
	}
	handler := createTestHandler(mockDB)

	req := httptest.NewRequest("GET", "/api/v1/data/p/stats", nil)
	rr := httptest.NewRecorder()

	handler.StatsEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["service"] != "data-service" {
		t.Errorf("Expected service 'data-service', got '%s'", response["service"])
	}

	dbStats, ok := response["database_stats"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected database_stats to be a map")
	}

	if dbStats["open_connections"] != float64(15) {
		t.Errorf("Expected open_connections 15, got %v", dbStats["open_connections"])
	}

	if dbStats["in_use"] != float64(5) {
		t.Errorf("Expected in_use 5, got %v", dbStats["in_use"])
	}

	if dbStats["idle"] != float64(10) {
		t.Errorf("Expected idle 10, got %v", dbStats["idle"])
	}
}

// TestHealthCheck_ContentType tests that Content-Type is always set
func TestHealthCheck_ContentType(t *testing.T) {
	testCases := []struct {
		name        string
		healthCheck func() error
	}{
		{
			name:        "Healthy",
			healthCheck: func() error { return nil },
		},
		{
			name:        "Unhealthy",
			healthCheck: func() error { return errors.New("error") },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := &MockDatabaseHandler{
				HealthCheckFunc: tc.healthCheck,
			}
			handler := createTestHandler(mockDB)

			req := httptest.NewRequest("GET", "/api/v1/data/p/health", nil)
			rr := httptest.NewRecorder()

			handler.HealthCheck(rr, req)

			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}
		})
	}
}
