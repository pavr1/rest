package sessionmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gateway-service/pkg/models"
	sharedHttp "shared/http"

	"github.com/sirupsen/logrus"
)

// SessionManager handles communication with the session service
type SessionManager struct {
	baseURL string
	client  *http.Client
	logger  *logrus.Logger
}

// NewSessionManager creates a new session manager
func NewSessionManager(sessionServiceURL string, logger *logrus.Logger) *SessionManager {
	return &SessionManager{
		baseURL: sessionServiceURL + "/api/v1/sessions",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// makeRequest makes a request to the session service with gateway headers
func (sm *SessionManager) makeRequest(method, path string, body io.Reader, requestID string) (*http.Response, error) {
	httpReq, err := http.NewRequest(method, sm.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Gateway-Service", "barrest-gateway")
	httpReq.Header.Set("X-Gateway-Session-Managed", "true")

	if requestID != "" {
		httpReq.Header.Set("X-Request-ID", requestID)
	}

	return sm.client.Do(httpReq)
}

// ValidateSession validates a token against the session service
func (sm *SessionManager) ValidateSession(token string, requestID string) (*models.TokenValidationResponse, error) {
	if token == "" {
		return &models.TokenValidationResponse{
			Valid:   false,
			Message: "Token is required",
		}, nil
	}

	validationReq := models.TokenValidationRequest{
		Token: token,
	}

	reqBody, err := json.Marshal(validationReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := sm.makeRequest("POST", "/p/validate", bytes.NewBuffer(reqBody), requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var responseWrapper sharedHttp.Response
	if err := json.Unmarshal(body, &responseWrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var validationResp models.TokenValidationResponse
	if responseWrapper.Data != nil {
		if dataBytes, err := json.Marshal(responseWrapper.Data); err == nil {
			json.Unmarshal(dataBytes, &validationResp)
		}
	}

	return &validationResp, nil
}

// LogoutSession revokes a session by token
func (sm *SessionManager) LogoutSession(token string, requestID string) error {
	req := models.SessionLogoutRequest{
		Token: token,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := sm.makeRequest("POST", "/logout", bytes.NewBuffer(reqBody), requestID)
	if err != nil {
		return fmt.Errorf("failed to logout session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logout failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
