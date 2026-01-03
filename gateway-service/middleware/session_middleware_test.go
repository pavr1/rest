package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractSessionIdFromHeader_Bearer(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-session-id-123")

	sessionId := extractSessionIdFromHeader(req)

	if sessionId != "test-session-id-123" {
		t.Errorf("extractSessionIdFromHeader() = %s; want test-session-id-123", sessionId)
	}
}

func TestExtractSessionIdFromHeader_Empty(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	sessionId := extractSessionIdFromHeader(req)

	if sessionId != "" {
		t.Errorf("extractSessionIdFromHeader() = %s; want empty string", sessionId)
	}
}

func TestExtractSessionIdFromHeader_NoBearerPrefix(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic some-credentials")

	sessionId := extractSessionIdFromHeader(req)

	if sessionId != "" {
		t.Errorf("extractSessionIdFromHeader() = %s; want empty string for non-Bearer auth", sessionId)
	}
}

func TestExtractSessionIdFromHeader_BearerOnly(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer ")

	sessionId := extractSessionIdFromHeader(req)

	if sessionId != "" {
		t.Errorf("extractSessionIdFromHeader() = %s; want empty string for 'Bearer ' only", sessionId)
	}
}

func TestNewSessionMiddleware(t *testing.T) {
	sm := NewSessionMiddleware(nil, nil)

	if sm == nil {
		t.Error("NewSessionMiddleware() returned nil")
	}
}

func TestSessionMiddleware_WriteErrorResponse(t *testing.T) {
	sm := &SessionMiddleware{}
	w := httptest.NewRecorder()

	sm.writeErrorResponse(w, http.StatusUnauthorized, "test_error", "Test message")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d; want %d", w.Code, http.StatusUnauthorized)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s; want application/json", contentType)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "test_error" {
		t.Errorf("error = %s; want test_error", response["error"])
	}

	if response["message"] != "Test message" {
		t.Errorf("message = %s; want Test message", response["message"])
	}

	if response["service"] != "gateway" {
		t.Errorf("service = %s; want gateway", response["service"])
	}

	if response["timestamp"] == nil {
		t.Error("timestamp should not be nil")
	}
}

func TestSessionMiddleware_ValidateSession_MissingAuth(t *testing.T) {
	sm := NewSessionMiddleware(nil, nil)

	nextCalled := false
	handler := sm.ValidateSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d; want %d", w.Code, http.StatusUnauthorized)
	}

	if nextCalled {
		t.Error("next handler should not be called when auth is missing")
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "missing_session" {
		t.Errorf("error = %s; want missing_session", response["error"])
	}
}
