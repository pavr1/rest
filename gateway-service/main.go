package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gateway-service/middleware"
	sessionmanager "gateway-service/middleware/session-manager"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	sharedConfig "shared/config"
	sharedLogger "shared/logger"
	sharedMiddlewares "shared/middlewares"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	// Gateway health check
	v1.HandleFunc("/gateway/p/health", createHealthHandler(sessionServiceUrl, dataServiceUrl, logger)).Methods("GET")

	// ==== PUBLIC ENDPOINTS (no authentication) ====

	// Session service - public endpoints
	sessionRouter := api.PathPrefix("/v1/sessions").Subrouter()
	sessionRouter.HandleFunc("/p/login", createProxyHandler(sessionServiceUrl, "/api/v1/sessions/p/login", logger)).Methods("POST")
	sessionRouter.HandleFunc("/p/validate", createProxyHandler(sessionServiceUrl, "/api/v1/sessions/p/validate", logger)).Methods("POST")
	sessionRouter.HandleFunc("/p/health", createProxyHandler(sessionServiceUrl, "/api/v1/sessions/p/health", logger)).Methods("GET")

	// Protected session endpoints
	sessionRouter.HandleFunc("/logout", createProxyHandler(sessionServiceUrl, "/api/v1/sessions/logout", logger)).Methods("POST")

	// Public health endpoints
	api.HandleFunc("/v1/data/p/health", createProxyHandler(dataServiceUrl, "/api/v1/data/p/health", logger)).Methods("GET")

	// ==== PROTECTED ENDPOINTS (require authentication) ====

	// Data service - protected endpoints
	dataRouter := api.PathPrefix("/v1/data").Subrouter()
	dataRouter.Use(sessionMiddleware.ValidateSession)
	dataRouter.PathPrefix("").HandlerFunc(createProxyHandler(dataServiceUrl, "/api/v1/data", logger))

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
		port = "8082"
	}

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

	log.Fatal(http.ListenAndServe(":"+port, r))
}

// createProxyHandler creates a reverse proxy handler for a specific service
func createProxyHandler(targetURL, stripPrefix string, logger *logrus.Logger) http.HandlerFunc {
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

func createHealthHandler(sessionServiceUrl, dataServiceUrl string, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionHealthy := checkServiceHealth(sessionServiceUrl+"/api/v1/sessions/p/health", logger)
		dataHealthy := checkServiceHealth(dataServiceUrl+"/api/v1/data/p/health", logger)

		status := "healthy"
		httpStatus := http.StatusOK

		if !sessionHealthy || !dataHealthy {
			status = "degraded"
			httpStatus = http.StatusServiceUnavailable
		}

		response := map[string]interface{}{
			"status":  status,
			"service": "gateway-service",
			"time":    time.Now(),
			"services": map[string]string{
				"gateway-service": "healthy",
				"session-service": boolToHealth(sessionHealthy),
				"data-service":    boolToHealth(dataHealthy),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(response)
	}
}

func checkServiceHealth(healthURL string, logger *logrus.Logger) bool {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", healthURL, nil)
	if err != nil {
		return false
	}

	req.Header.Set("X-Gateway-Service", "barrest-gateway")
	req.Header.Set("X-User-ID", "system")
	req.Header.Set("X-User-Role", "admin")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func boolToHealth(healthy bool) string {
	if healthy {
		return "healthy"
	}
	return "unhealthy"
}
