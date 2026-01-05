package middlewares

import (
	"net/http"
)

// GatewayHeaders contains the headers set by the gateway
type GatewayHeaders struct {
	Service        string
	SessionManaged bool
	UserID         string
	Username       string
	UserRole       string
	RequestID      string
}

// ExtractGatewayHeaders extracts gateway headers from the request
func ExtractGatewayHeaders(r *http.Request) GatewayHeaders {
	return GatewayHeaders{
		Service:        r.Header.Get("X-Gateway-Service"),
		SessionManaged: r.Header.Get("X-Gateway-Session-Managed") == "true",
		UserID:         r.Header.Get("X-User-ID"),
		Username:       r.Header.Get("X-Username"),
		UserRole:       r.Header.Get("X-User-Role"),
		RequestID:      r.Header.Get("X-Request-ID"),
	}
}

// GatewayMiddleware validates that requests come through the gateway
func GatewayMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now, just pass through
		// In production, you might want to validate that requests come from the gateway
		next.ServeHTTP(w, r)
	})
}

