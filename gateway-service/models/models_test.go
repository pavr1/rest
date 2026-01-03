package models

import (
	"encoding/json"
	"testing"
)

func TestSessionValidationRequest_JSON(t *testing.T) {
	req := SessionValidationRequest{
		SessionID: "test-session-123",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded SessionValidationRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.SessionID != req.SessionID {
		t.Errorf("SessionID = %s; want %s", decoded.SessionID, req.SessionID)
	}
}

func TestSessionValidationResponse_JSON(t *testing.T) {
	resp := SessionValidationResponse{
		Valid:       true,
		SessionID:   "session-123",
		Message:     "Session valid",
		UserID:      "user-456",
		Username:    "testuser",
		RoleName:    "admin",
		FullName:    "Test User",
		Permissions: []string{"read", "write"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded SessionValidationResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Valid != resp.Valid {
		t.Errorf("Valid = %v; want %v", decoded.Valid, resp.Valid)
	}
	if decoded.Username != resp.Username {
		t.Errorf("Username = %s; want %s", decoded.Username, resp.Username)
	}
	if len(decoded.Permissions) != len(resp.Permissions) {
		t.Errorf("Permissions length = %d; want %d", len(decoded.Permissions), len(resp.Permissions))
	}
}

func TestSessionCreateRequest_JSON(t *testing.T) {
	req := SessionCreateRequest{
		Username: "testuser",
		Password: "password123",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded SessionCreateRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Username != req.Username {
		t.Errorf("Username = %s; want %s", decoded.Username, req.Username)
	}
}

func TestSessionLogoutRequest_JSON(t *testing.T) {
	req := SessionLogoutRequest{
		SessionID: "session-to-logout",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded SessionLogoutRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.SessionID != req.SessionID {
		t.Errorf("SessionID = %s; want %s", decoded.SessionID, req.SessionID)
	}
}
