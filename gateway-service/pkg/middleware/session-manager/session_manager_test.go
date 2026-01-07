package sessionmanager

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sharedHttp "shared/http"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager("http://localhost:8087", nil)

	if sm == nil {
		t.Fatal("NewSessionManager() returned nil")
	}

	expectedBaseURL := "http://localhost:8087/api/v1/sessions"
	if sm.baseURL != expectedBaseURL {
		t.Errorf("baseURL = %s; want %s", sm.baseURL, expectedBaseURL)
	}

	if sm.client == nil {
		t.Error("client should not be nil")
	}

	if sm.client.Timeout.Seconds() != 10 {
		t.Errorf("client timeout = %v; want 10s", sm.client.Timeout)
	}
}

func TestSessionManager_ValidateSession_EmptySessionId(t *testing.T) {
	sm := NewSessionManager("http://localhost:8087", nil)

	resp, err := sm.ValidateSession("", "request-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Valid {
		t.Error("Valid should be false for empty session ID")
	}

	if resp.Message != "Session ID is required" {
		t.Errorf("Message = %s; want 'Session ID is required'", resp.Message)
	}
}

func TestSessionManager_ValidateSession_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("method = %s; want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/sessions/p/validate" {
			t.Errorf("path = %s; want /api/v1/sessions/p/validate", r.URL.Path)
		}
		if r.Header.Get("X-Gateway-Service") != "barrest-gateway" {
			t.Errorf("X-Gateway-Service = %s; want barrest-gateway", r.Header.Get("X-Gateway-Service"))
		}
		if r.Header.Get("X-Request-ID") != "req-123" {
			t.Errorf("X-Request-ID = %s; want req-123", r.Header.Get("X-Request-ID"))
		}

		// Return valid response
		response := sharedHttp.Response{
			Code:    200,
			Message: "Success",
			Data: map[string]interface{}{
				"valid":     true,
				"user_id":   "user-456",
				"username":  "testuser",
				"role_name": "admin",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	sm := NewSessionManager(server.URL, nil)

	resp, err := sm.ValidateSession("session-123", "req-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Valid {
		t.Error("Valid should be true")
	}

	if resp.UserID != "user-456" {
		t.Errorf("UserID = %s; want user-456", resp.UserID)
	}

	if resp.Username != "testuser" {
		t.Errorf("Username = %s; want testuser", resp.Username)
	}

	if resp.RoleName != "admin" {
		t.Errorf("RoleName = %s; want admin", resp.RoleName)
	}
}

func TestSessionManager_ValidateSession_Invalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := sharedHttp.Response{
			Code:    200,
			Message: "Success",
			Data: map[string]interface{}{
				"valid":   false,
				"message": "Session expired",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	sm := NewSessionManager(server.URL, nil)

	resp, err := sm.ValidateSession("expired-session", "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Valid {
		t.Error("Valid should be false for expired session")
	}

	if resp.Message != "Session expired" {
		t.Errorf("Message = %s; want 'Session expired'", resp.Message)
	}
}

func TestSessionManager_ValidateSession_ServiceDown(t *testing.T) {
	sm := NewSessionManager("http://localhost:19999", nil)

	_, err := sm.ValidateSession("session-123", "")

	if err == nil {
		t.Error("expected error when service is down")
	}
}

func TestSessionManager_LogoutSession_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s; want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/sessions/logout" {
			t.Errorf("path = %s; want /api/v1/sessions/logout", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sm := NewSessionManager(server.URL, nil)

	err := sm.LogoutSession("session-123", "req-456")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSessionManager_LogoutSession_Failed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid session"))
	}))
	defer server.Close()

	sm := NewSessionManager(server.URL, nil)

	err := sm.LogoutSession("invalid-session", "")

	if err == nil {
		t.Error("expected error for failed logout")
	}
}

func TestSessionManager_LogoutSession_ServiceDown(t *testing.T) {
	sm := NewSessionManager("http://localhost:19999", nil)

	err := sm.LogoutSession("session-123", "")

	if err == nil {
		t.Error("expected error when service is down")
	}
}
