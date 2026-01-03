package main

import (
	"net/http"
	"session-service/entities/sessions/handlers"
	sharedConfig "shared/config"
	httpresponse "shared/http-response"

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
	httpresponse.SendSuccess(w, http.StatusOK, "Session service healthy", map[string]string{
		"status":  "healthy",
		"service": "session-service",
	})
}
