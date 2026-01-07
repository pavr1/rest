package main

import (
	"context"
	"gateway-service/pkg/handlers"
	"net/http"
	"os"
	"os/signal"
	sharedConfig "shared/config"
	sharedLogger "shared/logger"
	"syscall"
	"time"
)

func main() {
	logger := sharedLogger.SetupLogger(sharedLogger.SERVICE_GATEWAY_SERVICE, "INFO")
	logger.Info("🌐 Gateway service starting")

	// Load configuration from data service
	configLoader := sharedConfig.NewConfigLoader(sharedConfig.DATA_SERVICE_URL)
	config, err := configLoader.LoadConfig("gateway", logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration from data service")
	}

	// Create HTTP handler and setup routes
	httpHandler := handlers.NewHTTPHandler(config, logger)
	router := httpHandler.SetupRoutes()

	// Start server
	port := config.GetString("SERVER_PORT")
	if port == "" {
		logger.Fatal("SERVER_PORT is not set")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		//pvillalobos this should be configurable
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("🚀 Gateway Service starting on :" + port)
		logger.Info("📡 Gateway Health: http://localhost:" + port + "/api/v1/gateway/p/health")
		logger.Info("")
		logger.Info("🔓 Public endpoints:")
		logger.Info("   GET  /api/v1/gateway/p/health       - Gateway health (checks business layer)")
		logger.Info("   POST /api/v1/sessions/p/login       - Login")
		logger.Info("   POST /api/v1/sessions/p/validate    - Validate session")
		logger.Info("   GET  /api/v1/sessions/p/health      - Session service health")
		logger.Info("")
		logger.Info("🔒 Protected endpoints (require Authorization header):")
		logger.Info("   POST /api/v1/sessions/logout        - Logout")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Fatal("Shutdown failed")
	}
	logger.Info("Gateway Service stopped")
}
