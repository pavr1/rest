package main

import (
	"net/http"
	"session-service/pkg/entities/sessions/handlers"
	sharedConfig "shared/config"
	sharedHealth "shared/health"
	httpresponse "shared/http-response"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	sessionsHandler *handlers.HTTPHandler
	healthMonitor   *sharedHealth.HealthMonitor
	logger          *logrus.Logger
}

func NewMainHTTPHandler(cfg *sharedConfig.Config, logger *logrus.Logger, healthMonitor *sharedHealth.HealthMonitor) (*MainHTTPHandler, error) {
	jwtHandler := handlers.NewJWTHandler(cfg.GetString("JWT_SECRET"), cfg.GetDuration("JWT_EXPIRATION_TIME"), logger)
	dbHandler, err := handlers.NewDBHandler(cfg, jwtHandler, logger)
	if err != nil {
		return nil, err
	}
	sessionsHandler := handlers.NewHTTPHandler(dbHandler, logger)
	return &MainHTTPHandler{
		sessionsHandler: sessionsHandler,
		healthMonitor:   healthMonitor,
		logger:          logger,
	}, nil
}

func (h *MainHTTPHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/sessions/p/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/api/v1/sessions/p/login", h.sessionsHandler.CreateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/p/validate", h.sessionsHandler.ValidateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/logout", h.sessionsHandler.LogoutSession).Methods("POST")
}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check if database connection is healthy
	if !h.sessionsHandler.GetDBHandler().IsHealthy() {
		httpresponse.SendError(w, http.StatusServiceUnavailable, "Database connection is not healthy", nil)
		return
	}

	// Use cached health state from background monitor
	if !h.healthMonitor.IsServiceHealthy("data-service") {
		httpresponse.SendError(w, http.StatusServiceUnavailable, "Data-service is not healthy", nil)
		return
	}

	status := h.healthMonitor.GetServiceStatus("data-service")
	httpresponse.SendSuccess(w, http.StatusOK, "Session service healthy", map[string]interface{}{
		"status":     "healthy",
		"service":    "session-service",
		"last_check": status.LastCheck,
	})
}
