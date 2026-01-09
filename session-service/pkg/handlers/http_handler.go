package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	sharedConfig "shared/config"
	sharedHttp "shared/http"

	sessionHandlers "session-service/pkg/entities/sessions/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	sessionsDBHandler   *sessionHandlers.DBHandler
	sessionsHandler     *sessionHandlers.HTTPHandler
	httpHealthMonitor   *sharedHttp.HTTPHealthMonitor
	cancelHealthMonitor context.CancelFunc
	logger              *logrus.Logger
}

func NewHTTPHandler(cfg *sharedConfig.Config, logger *logrus.Logger) (*MainHTTPHandler, error) {
	// Create JWT handler
	jwtHandler := sessionHandlers.NewJWTHandler(cfg.GetString("JWT_SECRET"), cfg.GetDuration("JWT_EXPIRATION_TIME"), logger)

	// Create sessions DB handler (creates its own DB connection)
	sessionsDBHandler, err := sessionHandlers.NewDBHandler(cfg, jwtHandler, logger)
	if err != nil {
		return nil, err
	}

	// Create sessions HTTP handler
	sessionsHTTPHandler := sessionHandlers.NewHTTPHandler(sessionsDBHandler, logger)

	// Create cancellable context for health monitor
	ctx, cancel := context.WithCancel(context.Background())

	//pvillalobos this should be configurable
	// Create HTTP health monitor for data-service
	httpHealthMonitor, err := sharedHttp.NewHealthMonitor(logger, 1*time.Second)
	if err != nil {
		cancel()
		return nil, err
	}
	httpHealthMonitor.AddService("data-service", sharedConfig.DATA_SERVICE_URL+"/api/v1/data/p/health")
	httpHealthMonitor.Start(ctx)

	return &MainHTTPHandler{
		sessionsDBHandler:   sessionsDBHandler,
		sessionsHandler:     sessionsHTTPHandler,
		httpHealthMonitor:   httpHealthMonitor,
		cancelHealthMonitor: cancel,
		logger:              logger,
	}, nil
}

func (h *MainHTTPHandler) CloseDB() error {
	// Stop health monitor
	if h.cancelHealthMonitor != nil {
		h.cancelHealthMonitor()
	}

	err := h.sessionsDBHandler.Close()
	if err != nil {
		h.logger.WithError(err).Error("Failed to close database")
		return err
	}
	return nil
}

func (h *MainHTTPHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/sessions/p/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/api/v1/sessions/p/login", h.sessionsHandler.CreateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/p/validate", h.sessionsHandler.ValidateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/logout", h.sessionsHandler.LogoutSession).Methods("POST")
}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "session-service",
		"timestamp": time.Now(),
	}

	// Check cached health state from background monitor
	healthStatus := h.httpHealthMonitor.GetHealthStatus()
	if !healthStatus.IsHealthy {
		response["status"] = "unhealthy"
		response["message"] = "Dependent services are not healthy"
		response["services"] = healthStatus.Services

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response["status"] = "healthy"
	response["message"] = "Session service is healthy"
	response["services"] = healthStatus.Services

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
