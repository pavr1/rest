package middleware

import (
	"encoding/json"
	sessionmanager "gateway-service/pkg/middleware/session-manager"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// SessionMiddleware handles session validation for protected routes
type SessionMiddleware struct {
	sessionManager *sessionmanager.SessionManager
	logger         *logrus.Logger
}

// NewSessionMiddleware creates a new session middleware
func NewSessionMiddleware(sessionManager *sessionmanager.SessionManager, logger *logrus.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// ValidateSession middleware validates the token against the session service
func (sm *SessionMiddleware) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		token := extractTokenFromHeader(r)
		if token == "" {
			sm.writeErrorResponse(w, http.StatusUnauthorized, "missing_token", "Token is required")
			return
		}

		requestID := r.Header.Get("X-Request-ID")

		// Validate token with session service
		validation, err := sm.sessionManager.ValidateSession(token, requestID)
		if err != nil {
			sm.logger.WithError(err).Error("Session validation error")
			sm.writeErrorResponse(w, http.StatusInternalServerError, "validation_error", "Failed to validate session")
			return
		}

		if !validation.Valid {
			sm.writeErrorResponse(w, http.StatusUnauthorized, "invalid_session", validation.Message)
			return
		}

		// Add user context to request headers for backend services
		r.Header.Set("X-User-ID", validation.StaffID)
		r.Header.Set("X-Username", validation.Username)
		r.Header.Set("X-User-Role", validation.Role)

		if len(validation.Permissions) > 0 {
			r.Header.Set("X-User-Permissions", strings.Join(validation.Permissions, ","))
		}

		next.ServeHTTP(w, r)
	})
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const bearerPrefix = "Bearer "
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return authHeader[len(bearerPrefix):]
	}

	return ""
}

func (sm *SessionMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"error":     errorCode,
		"message":   message,
		"timestamp": time.Now(),
		"service":   "gateway",
	}

	json.NewEncoder(w).Encode(response)
}
