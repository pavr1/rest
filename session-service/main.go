package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"session-service/pkg/entities/sessions/handlers"
	"syscall"
	"time"

	sharedConfig "shared/config"
	sharedHealth "shared/health"
	httpresponse "shared/http-response"
	sharedLogger "shared/logger"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	sessionsHandler *handlers.HTTPHandler
	logger          *logrus.Logger
	healthMonitor   *sharedHealth.HealthMonitor
}

const (
	// pvillalobos this should be configurable
	HealthCheckInterval = 1 * time.Second
)

func main() {
	logger := sharedLogger.SetupLogger(sharedLogger.SERVICE_SESSION_SERVICE, "INFO")
	logger.Info("🔐 Starting Bar-Restaurant Session Service")

	configLoader := sharedConfig.NewConfigLoader(sharedConfig.DATA_SERVICE_URL)
	config, err := configLoader.LoadConfig("session", logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	logger.WithFields(logrus.Fields{
		"port": config.GetString("SERVER_PORT"),
		"host": config.GetString("SERVER_HOST"),
	}).Info("Configuration loaded")

	// Start health monitor for database in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	healthMonitor := sharedHealth.NewHealthMonitor(logger, HealthCheckInterval)
	healthMonitor.AddService("data-service", sharedConfig.DATA_SERVICE_URL+"/api/v1/data/p/health")
	go healthMonitor.Start(ctx)

	mainHandler, dbHandler, err := newMainHTTPHandler(config, logger, healthMonitor)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create HTTP handler")
	}
	defer dbHandler.Close()

	router := mux.NewRouter()
	mainHandler.SetupRoutes(router)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GetString("SERVER_HOST"), config.GetString("SERVER_PORT")),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.WithField("addr", server.Addr).Info("HTTP server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down...")
	cancel() // Stop health monitoring

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Fatal("Shutdown failed")
	}
	logger.Info("Session Service stopped")
}

func newMainHTTPHandler(cfg *sharedConfig.Config, logger *logrus.Logger, healthMonitor *sharedHealth.HealthMonitor) (*MainHTTPHandler, *handlers.DBHandler, error) {
	jwtHandler := handlers.NewJWTHandler(cfg.GetString("JWT_SECRET"), cfg.GetDuration("JWT_EXPIRATION_TIME"), logger)
	dbHandler, err := handlers.NewDBHandler(cfg, jwtHandler, logger)
	if err != nil {
		return nil, nil, err
	}
	sessionsHandler := handlers.NewHTTPHandler(dbHandler, logger)
	return &MainHTTPHandler{
		sessionsHandler: sessionsHandler,
		logger:          logger,
		healthMonitor:   healthMonitor,
	}, dbHandler, nil
}

func (h *MainHTTPHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/sessions/p/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/api/v1/sessions/p/login", h.sessionsHandler.CreateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/p/validate", h.sessionsHandler.ValidateSession).Methods("POST")
	router.HandleFunc("/api/v1/sessions/logout", h.sessionsHandler.LogoutSession).Methods("POST")
}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check shared database health state (updated by main.go background goroutine)
	if !h.healthMonitor.AreAllServicesHealthy() {
		httpresponse.SendError(w, http.StatusServiceUnavailable, "Database is not healthy", nil)
		return
	}

	httpresponse.SendSuccess(w, http.StatusOK, "Session service healthy", map[string]interface{}{
		"status":  "healthy",
		"service": "session-service",
	})
}
