package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCreateSessionBadRequest(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create handler with nil dbHandler (will fail on actual login)
	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	// Test with invalid JSON
	req := httptest.NewRequest("POST", "/api/v1/sessions/p/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateSession(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestCreateSessionMissingCredentials(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	// Test with empty username/password
	body := `{"username":"","password":""}`
	req := httptest.NewRequest("POST", "/api/v1/sessions/p/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateSession(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusBadRequest)
	}

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["message"] != "Username and password are required" {
		t.Errorf("message = %q, want %q", response["message"], "Username and password are required")
	}
}

func TestCreateSessionMissingUsername(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	body := `{"username":"","password":"somepassword"}`
	req := httptest.NewRequest("POST", "/api/v1/sessions/p/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateSession(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestCreateSessionMissingPassword(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	body := `{"username":"admin","password":""}`
	req := httptest.NewRequest("POST", "/api/v1/sessions/p/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateSession(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestValidateSessionBadRequest(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/v1/sessions/p/validate", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ValidateSession(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestLogoutSessionBadRequest(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/v1/sessions/logout", bytes.NewBufferString("bad json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.LogoutSession(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestLogoutSessionMissingSessionID(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	body := `{"session_id":""}`
	req := httptest.NewRequest("POST", "/api/v1/sessions/logout", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.LogoutSession(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rr.Code, http.StatusBadRequest)
	}

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["message"] != "Session ID required" {
		t.Errorf("message = %q, want %q", response["message"], "Session ID required")
	}
}

func TestContentTypeJSON(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	handler := &HTTPHandler{
		dbHandler: nil,
		logger:    logger,
	}

	body := `{"username":"","password":""}`
	req := httptest.NewRequest("POST", "/api/v1/sessions/p/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateSession(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}
}
