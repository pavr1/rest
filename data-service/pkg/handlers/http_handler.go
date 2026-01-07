package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	sharedDb "shared/db"

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
	db     *sharedDb.DbHandler
	config *sharedDb.Config
	logger *logrus.Logger
}

// NewHandler creates a new HTTP handler
func NewHTTPHandler(db *sharedDb.DbHandler, config *sharedDb.Config, logger *logrus.Logger) (*HTTPHandler, error) {
	// repository, err := settings.NewRepository(db)
	// if err != nil {
	// 	return nil, err
	// }
	// settingsHandler := settingsHTTP.NewHTTPHandler(repository, logger)

	return &HTTPHandler{
		//settingsHandler: settingsHandler,
		db:     db,
		config: config,
		logger: logger,
	}, nil
}

// SetupRoutes configures all HTTP routes
func (h *HTTPHandler) SetupRoutes(router *mux.Router) {
	//Root endpoint
	router.HandleFunc("/", h.RootHandler).Methods("GET")

	//Public endpoints
	router.HandleFunc("/api/v1/data/p/health", h.HealthCheck).Methods("GET")

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
	if !h.db.IsConnected() {
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
