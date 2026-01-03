package main

import (
	"net/http"
	"session-service/entities/sessions/handlers"
	sharedConfig "shared/config"
	httpresponse "shared/http-response"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	sessionsHandler *handlers.HTTPHandler
	logger          *logrus.Logger
}

func NewMainHTTPHandler(cfg *sharedConfig.Config, logger *logrus.Logger) (*MainHTTPHandler, error) {
	jwtHandler := handlers.NewJWTHandler(cfg.GetString("JWT_SECRET"), cfg.GetDuration("JWT_EXPIRATION_TIME"), logger)
	dbHandler, err := handlers.NewDBHandler(cfg, jwtHandler, logger)
	if err != nil {
		return nil, err
	}
	sessionsHandler := handlers.NewHTTPHandler(dbHandler, logger)
	return &MainHTTPHandler{sessionsHandler: sessionsHandler, logger: logger}, nil
}

func (h *MainHTTPHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/sessions/p/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/api/v1/sessions/p/login", h.sessionsHandler.CreateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/p/validate", h.sessionsHandler.ValidateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/logout", h.sessionsHandler.LogoutSession).Methods("POST")
}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check data-service health first
	dataServiceHealthy := h.checkDataServiceHealth(r)

	if !dataServiceHealthy {
		httpresponse.SendError(w, http.StatusServiceUnavailable, "Data-service is not healthy", nil)
		return
	}

	httpresponse.SendSuccess(w, http.StatusOK, "Session service healthy", map[string]string{
		"status":  "healthy",
		"service": "session-service",
	})
}

// checkDataServiceHealth checks if the data-service is healthy
func (h *MainHTTPHandler) checkDataServiceHealth(r *http.Request) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Direct call to data service (internal service communication)
	req, err := http.NewRequest("GET", sharedConfig.DATA_SERVICE_URL+"/api/v1/data/p/health", nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create data service health check request")
		return false
	}

	// Add headers for internal service communication
	req.Header.Set("X-Gateway-Service", "session-service")
	req.Header.Set("X-User-ID", "system")
	req.Header.Set("X-User-Role", "admin")

	// Forward the existing X-Request-ID from the current request
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	resp, err := client.Do(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to connect to data-service")
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
