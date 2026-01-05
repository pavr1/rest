package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gateway-service/pkg/middleware"
	sessionmanager "gateway-service/pkg/middleware/session-manager"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	sharedConfig "shared/config"
	sharedHealth "shared/health"
	sharedLogger "shared/logger"
	sharedMiddlewares "shared/middlewares"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	// HealthCheckInterval is how often to ping dependencies
	HealthCheckInterval = 10 * time.Second
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

	// Service URLs
	sessionServiceUrl := config.GetString("SESSION_SERVICE_URL")
	dataServiceUrl := config.GetString("DATA_SERVICE_URL")

	logger.WithFields(logrus.Fields{
		"session_service": sessionServiceUrl,
		"data_service":    dataServiceUrl,
	}).Info("Configuration loaded")

	// Start health monitor (background goroutine)
	ctx, cancel := context.WithCancel(context.Background())
	healthMonitor := sharedHealth.NewHealthMonitor(logger, HealthCheckInterval)
	healthMonitor.AddService("data-service", dataServiceUrl+"/api/v1/data/p/health")
	healthMonitor.AddService("session-service", sessionServiceUrl+"/api/v1/sessions/p/health")
	// Future services - add here as they are created:
	// healthMonitor.AddService("menu-service", menuServiceUrl+"/api/v1/menu/p/health")
	// healthMonitor.AddService("orders-service", ordersServiceUrl+"/api/v1/orders/p/health")
	go healthMonitor.Start(ctx)

	// Create session manager for authentication
	sessionManager := sessionmanager.NewSessionManager(sessionServiceUrl, logger)
	sessionMiddleware := middleware.NewSessionMiddleware(sessionManager, logger)

	r := mux.NewRouter()

	// Apply global middleware
	r.Use(sharedMiddlewares.RequestIDMiddleware)
	r.Use(sharedMiddlewares.GatewayMiddleware)

	// CORS middleware
	corsMiddleware := middleware.NewCORSMiddleware(logger)
	r.Use(corsMiddleware.HandleCORS)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()

	// Gateway health check (uses cached health from monitor)
	v1.HandleFunc("/gateway/p/health", createHealthHandler(healthMonitor)).Methods("GET")

	// ==== PUBLIC ENDPOINTS (no authentication) ====

	// Session service - public endpoints
	api.HandleFunc("/v1/sessions/p/login", createProxyHandler(sessionServiceUrl, logger)).Methods("POST")
	api.HandleFunc("/v1/sessions/p/validate", createProxyHandler(sessionServiceUrl, logger)).Methods("POST")
	api.HandleFunc("/v1/sessions/p/health", createProxyHandler(sessionServiceUrl, logger)).Methods("GET")

	// // Public health endpoints
	// api.HandleFunc("/v1/data/p/health", createProxyHandler(dataServiceUrl, logger)).Methods("GET")

	// // Data service - public settings endpoint (for service-to-service config loading)
	// api.HandleFunc("/v1/data/settings/by-service", createProxyHandler(dataServiceUrl, logger)).Methods("POST")

	// ==== PROTECTED SESSION ENDPOINTS (require authentication) ====
	protectedSessionRouter := api.PathPrefix("/v1/sessions").Subrouter()
	protectedSessionRouter.Use(sessionMiddleware.ValidateSession)
	protectedSessionRouter.HandleFunc("/logout", createProxyHandler(sessionServiceUrl, logger)).Methods("POST")

	// ==== PROTECTED ENDPOINTS (require authentication) ====

	// // Data service - protected endpoints
	// dataRouter := api.PathPrefix("/v1/data").Subrouter()
	// dataRouter.Use(sessionMiddleware.ValidateSession)
	// dataRouter.PathPrefix("").HandlerFunc(createProxyHandler(dataServiceUrl, logger))

	// Future service routes (to be added as services are created):
	// - /api/v1/menu/* → menu-service
	// - /api/v1/orders/* → orders-service
	// - /api/v1/inventory/* → inventory-service
	// - /api/v1/payments/* → payment-service
	// - /api/v1/invoices/* → invoice-service
	// - /api/v1/karaoke/* → karaoke-service
	// - /api/v1/promotions/* → promotion-service
	// - /api/v1/customers/* → customer-service

	// OPTIONS handling for CORS preflight
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Start server
	port := config.GetString("SERVER_PORT")
	if port == "" {
		logger.Fatal("SERVER_PORT is not set")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		//pvillalobos this should be configurable
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("🚀 Gateway Service starting on :" + port)
		logger.Info("📡 Health: http://localhost:" + port + "/api/v1/gateway/p/health")
		logger.Info("")
		logger.Info("🔓 Public endpoints:")
		logger.Info("   POST /api/v1/sessions/p/login")
		logger.Info("   POST /api/v1/sessions/p/validate")
		logger.Info("   GET  /api/v1/sessions/p/health")
		logger.Info("   GET  /api/v1/data/p/health")
		logger.Info("")
		logger.Info("🔒 Protected endpoints (require Authorization header):")
		logger.Info("   POST /api/v1/sessions/logout")
		logger.Info("   ALL  /api/v1/data/*")

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
	logger.Info("Gateway Service stopped")
}

// createProxyHandler creates a reverse proxy handler for a specific service
func createProxyHandler(targetURL string, logger *logrus.Logger) http.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		logger.Fatalf("Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		serviceName := "unknown-service"
		switch {
		case strings.Contains(r.URL.Path, "/sessions"):
			serviceName = "session-service"
		case strings.Contains(r.URL.Path, "/data"):
			serviceName = "data-service"
		case strings.Contains(r.URL.Path, "/orders"):
			serviceName = "orders-service"
		case strings.Contains(r.URL.Path, "/menu"):
			serviceName = "menu-service"
		case strings.Contains(r.URL.Path, "/inventory"):
			serviceName = "inventory-service"
		}

		logger.WithFields(logrus.Fields{
			"service": serviceName,
			"path":    r.URL.Path,
			"error":   err.Error(),
		}).Error("Proxy error - service unavailable")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":     "service_unavailable",
			"message":   fmt.Sprintf("The %s is currently unavailable", serviceName),
			"timestamp": time.Now(),
			"service":   serviceName,
		})
	}

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		requestID := req.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			req.Header.Set("X-Request-ID", requestID)
		}

		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Gateway-Service", "barrest-gateway")
		req.Header.Set("X-Gateway-Session-Managed", "true")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func createHealthHandler(healthMonitor *sharedHealth.HealthMonitor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use cached health state from background monitor
		allServices := healthMonitor.GetAllServicesStatus()

		gatewayHealthy := true
		httpStatus := http.StatusOK

		// Check if all services are healthy
		services := make(map[string]interface{})

		for name, svc := range allServices {
			if svc.Healthy {
				services[name] = "healthy"
			} else {
				services[name] = "unhealthy"
				// Gateway depends on all its services - if any is down, gateway is unhealthy
				gatewayHealthy = false
				httpStatus = http.StatusServiceUnavailable
			}
		}

		// Gateway status depends on its dependencies
		if gatewayHealthy {
			services["gateway-service"] = "healthy"
		} else {
			services["gateway-service"] = "unhealthy"
		}

		response := map[string]interface{}{
			"status":   map[bool]string{true: "healthy", false: "unhealthy"}[gatewayHealthy],
			"service":  "gateway-service",
			"time":     time.Now(),
			"services": services,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(response)
	}
}
