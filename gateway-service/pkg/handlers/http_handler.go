package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gateway-service/pkg/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	sharedConfig "shared/config"
	sharedHttp "shared/http"
	sharedMiddlewares "shared/middlewares"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	config            *sharedConfig.Config
	sessionServiceUrl string
	httpHealthMonitor *sharedHttp.HTTPHealthMonitor
	logger            *logrus.Logger
}

func NewHTTPHandler(
	config *sharedConfig.Config,
	sessionServiceUrl string,
	httpHealthMonitor *sharedHttp.HTTPHealthMonitor,
	logger *logrus.Logger,
) *HTTPHandler {
	return &HTTPHandler{
		config:            config,
		sessionServiceUrl: sessionServiceUrl,
		httpHealthMonitor: httpHealthMonitor,
		logger:            logger,
	}
}

// GatewayHealthCheck handles the gateway health check endpoint
// Uses cached health state from background health monitor
func (h *HTTPHandler) GatewayHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Get health status from monitor (cached, updated every 10s)
	healthStatus := h.httpHealthMonitor.GetHealthStatus()

	// Set status code based on overall health
	statusCode := http.StatusOK
	if !healthStatus.IsHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(healthStatus)
}

// CreateProxyHandler creates a reverse proxy handler for a specific service
func (h *HTTPHandler) CreateProxyHandler(targetURL string) http.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		h.logger.Fatalf("Invalid target URL: %v", err)
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

		h.logger.WithFields(logrus.Fields{
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

// SetupRoutes configures all gateway routes
func (h *HTTPHandler) SetupRoutes(sessionMiddleware *middleware.SessionMiddleware) *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(sharedMiddlewares.RequestIDMiddleware)
	r.Use(sharedMiddlewares.GatewayMiddleware)

	// CORS middleware
	corsMiddleware := middleware.NewCORSMiddleware(h.logger)
	r.Use(corsMiddleware.HandleCORS)

	api := r.PathPrefix("/api").Subrouter()

	// Gateway health endpoint (checks business layer services only)
	api.HandleFunc("/v1/gateway/p/health", h.GatewayHealthCheck).Methods("GET")

	// ==== PUBLIC ENDPOINTS (no authentication) ====

	// Session service - public endpoints
	api.HandleFunc("/v1/sessions/p/login", h.CreateProxyHandler(h.sessionServiceUrl)).Methods("POST")
	api.HandleFunc("/v1/sessions/p/validate", h.CreateProxyHandler(h.sessionServiceUrl)).Methods("POST")
	api.HandleFunc("/v1/sessions/p/health", h.CreateProxyHandler(h.sessionServiceUrl)).Methods("GET")

	// ==== PROTECTED SESSION ENDPOINTS (require authentication) ====
	protectedSessionRouter := api.PathPrefix("/v1/sessions").Subrouter()
	protectedSessionRouter.Use(sessionMiddleware.ValidateSession)
	protectedSessionRouter.HandleFunc("/logout", h.CreateProxyHandler(h.sessionServiceUrl)).Methods("POST")

	// OPTIONS handling for CORS preflight
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
