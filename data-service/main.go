package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"data-service/pkg/database"
	sharedConfig "shared/config"
	sharedLogger "shared/logger"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	// Setup HTTP server
	router := setupRouter(db, config, logger)

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

// setupRouter configures the HTTP routes
func setupRouter(db database.DatabaseHandler, config *database.Config, logger *logrus.Logger) *mux.Router {
	router := mux.NewRouter()

	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Bar-Restaurant Data Service is running"}`))
	}).Methods("GET")

	// Public health check endpoint
	router.HandleFunc("/api/v1/data/p/health", func(w http.ResponseWriter, r *http.Request) {
		healthCheck(w, db, config, logger)
	}).Methods("GET")

	// Stats endpoint
	router.HandleFunc("/api/v1/data/p/stats", func(w http.ResponseWriter, r *http.Request) {
		statsEndpoint(w, r, db, logger)
	}).Methods("GET")

	return router
}

// healthCheck handles the health check endpoint
func healthCheck(w http.ResponseWriter, db database.DatabaseHandler, config *database.Config, logger *logrus.Logger) {
	response := map[string]interface{}{
		"service":   "data-service",
		"timestamp": time.Now(),
	}

	if err := db.HealthCheck(); err != nil {
		logger.WithError(err).Error("Database ping check failed")
		response["status"] = "unhealthy"
		response["message"] = "Database ping check failed"
		response["error"] = err.Error()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response["status"] = "healthy"
	response["message"] = "Database ping check passed"
	response["database"] = map[string]interface{}{
		"host":   config.Host,
		"port":   config.Port,
		"dbname": config.DBName,
		"stats":  db.GetStats(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// statsEndpoint provides database connection statistics
func statsEndpoint(w http.ResponseWriter, _ *http.Request, db database.DatabaseHandler, _ *logrus.Logger) {
	stats := db.GetStats()

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
