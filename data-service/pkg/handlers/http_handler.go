package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HealthChecker interface for checking health state
type HealthChecker interface {
	IsHealthy() bool
}

// Handler is the main HTTP handler for data-service
type HTTPHandler struct {
	//settingsHandler *settingsHTTP.HTTPHandler
	db            IDBHandler
	config        *Config
	logger        *logrus.Logger
	healthChecker HealthChecker
}

// NewHandler creates a new HTTP handler
func NewHTTPHandler(db IDBHandler, config *Config, logger *logrus.Logger, healthChecker HealthChecker) (*HTTPHandler, error) {
	// repository, err := settings.NewRepository(db)
	// if err != nil {
	// 	return nil, err
	// }
	// settingsHandler := settingsHTTP.NewHTTPHandler(repository, logger)

	return &HTTPHandler{
		//settingsHandler: settingsHandler,
		db:            db,
		config:        config,
		logger:        logger,
		healthChecker: healthChecker,
	}, nil
}

// SetupRoutes configures all HTTP routes
func (h *HTTPHandler) SetupRoutes(router *mux.Router) {
	//Root endpoint
	router.HandleFunc("/", h.RootHandler).Methods("GET")

	//Public endpoints
	router.HandleFunc("/api/v1/data/p/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/api/v1/data/p/stats", h.StatsEndpoint).Methods("GET")

	//pvillalobos these settings must belong to settings service
	// //Settings endpoints
	// router.HandleFunc("/api/v1/data/settings/by-service", h.settingsHandler.GetSettingsByService).Methods("POST")
	// router.HandleFunc("/api/v1/data/settings/by-key", h.settingsHandler.GetSettingByKey).Methods("POST")
	// router.HandleFunc("/api/v1/data/settings", h.settingsHandler.CreateSetting).Methods("POST")
	// router.HandleFunc("/api/v1/data/settings", h.settingsHandler.UpdateSetting).Methods("PUT")
	// router.HandleFunc("/api/v1/data/settings", h.settingsHandler.DeleteSetting).Methods("DELETE")
}

// RootHandler handles the root endpoint
func (h *HTTPHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Bar-Restaurant Data Service is running"}`))
}

// HealthCheck handles the health check endpoint
// Returns cached health state that's updated by background ping loop every second
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "data-service",
		"timestamp": time.Now(),
	}

	// Check cached health state (updated by background health monitor in main.go)
	if !h.healthChecker.IsHealthy() {
		response["status"] = "unhealthy"
		response["message"] = "Database is not reachable"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response["status"] = "healthy"
	response["message"] = "Database is healthy"
	response["database"] = map[string]interface{}{
		"host":   h.config.Host,
		"port":   h.config.Port,
		"dbname": h.config.DBName,
		"stats":  h.db.GetStats(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// StatsEndpoint provides database connection statistics
func (h *HTTPHandler) StatsEndpoint(w http.ResponseWriter, r *http.Request) {
	stats := h.db.GetStats()

	response := map[string]interface{}{
		"service":   "data-service",
		"timestamp": time.Now(),
		"database_stats": map[string]interface{}{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
			"wait_count":       stats.WaitCount,
			"wait_duration":    stats.WaitDuration.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
