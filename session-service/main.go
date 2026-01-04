package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sharedConfig "shared/config"
	sharedHealth "shared/health"
	sharedLogger "shared/logger"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	// HealthCheckInterval is how often to ping dependencies
	HealthCheckInterval = 10 * time.Second
)

func main() {
	logger := sharedLogger.SetupLogger(sharedLogger.SERVICE_SESSION_SERVICE, "INFO")
	logger.Info("🔐 Starting Bar-Restaurant Session Service")

	configLoader := sharedConfig.NewConfigLoader(sharedConfig.DATA_SERVICE_URL)
	config, err := configLoader.LoadConfig("session", logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	if len(config.Values) == 0 {
		logger.Fatal("No configuration values loaded")
	}

	logger.WithFields(logrus.Fields{
		"port": config.GetString("SERVER_PORT"),
		"host": config.GetString("SERVER_HOST"),
	}).Info("Configuration loaded")

	// Start health monitor (background goroutine)
	ctx, cancel := context.WithCancel(context.Background())
	healthMonitor := sharedHealth.NewHealthMonitor(logger, HealthCheckInterval)
	healthMonitor.AddService("data-service", sharedConfig.DATA_SERVICE_URL+"/api/v1/data/p/health")
	go healthMonitor.Start(ctx)

	mainHandler, err := NewMainHTTPHandler(config, logger, healthMonitor)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create HTTP handler")
	}

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
	cancel() // Stop health monitor

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Fatal("Shutdown failed")
	}
	logger.Info("Session Service stopped")
}
