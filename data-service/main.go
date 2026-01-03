package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"data-service/pkg/database"
	httpHandler "data-service/pkg/http"
	sharedConfig "shared/config"
	sharedLogger "shared/logger"

	"github.com/gorilla/mux"
)

func main() {
	logger := sharedLogger.SetupLogger(sharedLogger.SERVICE_DATA_SERVICE, "INFO")

	config := database.DefaultConfig(logger)

	// Create database handler
	db := database.New(config, logger)

	// Connect to database
	fmt.Println("🍺 Connecting to Bar-Restaurant Data Service...")
	if err := db.Connect(); err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Perform initial health check
	if err := db.HealthCheck(); err != nil {
		logger.WithError(err).Fatal("Initial database health check failed")
	}

	fmt.Println("✅ Database connection established successfully")

	// Setup HTTP handler and router
	handler, err := httpHandler.NewHandler(db, config, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create HTTP handler")
	}
	router := mux.NewRouter()
	handler.SetupRoutes(router)

	// Get server configuration
	serverHost := sharedConfig.DATA_SERVICE_HOST
	serverPort := sharedConfig.DATA_SERVICE_PORT

	logger.WithField("DATA_SERVICE_HOST", serverHost).Info("🔍 DEBUG: DATA_SERVICE_HOST")
	logger.WithField("DATA_SERVICE_PORT", serverPort).Info("🔍 DEBUG: DATA_SERVICE_PORT")

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", serverHost, serverPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.WithField("port", serverPort).Info("Starting Data Service HTTP server")
		logger.WithField("port", serverPort).Info("🚀 Data Service HTTP server starting on :8086")
		logger.WithField("port", serverPort).Info("📡 Health endpoint available at: http://localhost:8086/api/v1/data/p/health")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Data Service...")

	// Gracefully shutdown with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Data Service exited gracefully")
}
