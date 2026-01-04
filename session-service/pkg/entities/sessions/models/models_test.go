package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSessionCreateRequest(t *testing.T) {
	req := SessionCreateRequest{
		Username: "testuser",
		Password: "testpass",
	}

	if req.Username != "testuser" {
		t.Errorf("Username = %q, want %q", req.Username, "testuser")
	}

	if req.Password != "testpass" {
		t.Errorf("Password = %q, want %q", req.Password, "testpass")
	}
}

func TestSessionCreateRequestJSON(t *testing.T) {
	jsonStr := `{"username":"admin","password":"secret123"}`

	var req SessionCreateRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if req.Username != "admin" {
		t.Errorf("Username = %q, want %q", req.Username, "admin")
	}

	if req.Password != "secret123" {
		t.Errorf("Password = %q, want %q", req.Password, "secret123")
	}
}

func TestSessionCreateResponse(t *testing.T) {
	staff := &Staff{
		ID:        "staff-123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Role:      "waiter",
	}

	resp := SessionCreateResponse{
		SessionID: "session-abc",
		Token:     "jwt-token-xyz",
		Message:   "Login successful",
		Staff:     staff,
	}

	if resp.SessionID != "session-abc" {
		t.Errorf("SessionID = %q, want %q", resp.SessionID, "session-abc")
	}

	if resp.Staff.Username != "testuser" {
		t.Errorf("Staff.Username = %q, want %q", resp.Staff.Username, "testuser")
	}
}

func TestSessionValidationResponse(t *testing.T) {
	resp := SessionValidationResponse{
		Valid:     true,
		SessionID: "session-123",
		Message:   "Valid session",
		StaffID:   "staff-456",
		Username:  "admin",
		Role:      "manager",
		FullName:  "Admin User",
	}

	if !resp.Valid {
		t.Error("Valid should be true")
	}

	if resp.Role != "manager" {
		t.Errorf("Role = %q, want %q", resp.Role, "manager")
	}
}

func TestSessionLogoutResponse(t *testing.T) {
	resp := SessionLogoutResponse{
		Success:   true,
		SessionID: "session-789",
		Message:   "Logged out successfully",
	}

	if !resp.Success {
		t.Error("Success should be true")
	}

	if resp.Message != "Logged out successfully" {
		t.Errorf("Message = %q, want %q", resp.Message, "Logged out successfully")
	}
}

func TestStaffPasswordHashHidden(t *testing.T) {
	staff := Staff{
		ID:           "staff-123",
		Username:     "testuser",
		PasswordHash: "super-secret-hash",
		FirstName:    "Test",
		LastName:     "User",
		Role:         "waiter",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(staff)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Password hash should NOT appear in JSON (json:"-" tag)
	if contains(jsonStr, "super-secret-hash") {
		t.Error("PasswordHash should not be in JSON output")
	}

	if contains(jsonStr, "password_hash") {
		t.Error("password_hash key should not be in JSON output")
	}

	// But other fields should be present
	if !contains(jsonStr, "testuser") {
		t.Error("Username should be in JSON output")
	}
}

func TestStaffOptionalFields(t *testing.T) {
	email := "test@example.com"
	loginTime := time.Now()

	staff := Staff{
		ID:          "staff-123",
		Username:    "testuser",
		Email:       &email,
		FirstName:   "Test",
		LastName:    "User",
		Role:        "chef",
		IsActive:    true,
		LastLoginAt: &loginTime,
	}

	if *staff.Email != email {
		t.Errorf("Email = %q, want %q", *staff.Email, email)
	}

	if staff.LastLoginAt == nil {
		t.Error("LastLoginAt should not be nil")
	}
}

func TestStaffNilOptionalFields(t *testing.T) {
	staff := Staff{
		ID:        "staff-123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Role:      "bartender",
		IsActive:  true,
	}

	if staff.Email != nil {
		t.Error("Email should be nil")
	}

	if staff.LastLoginAt != nil {
		t.Error("LastLoginAt should be nil")
	}

	// Should marshal without error
	_, err := json.Marshal(staff)
	if err != nil {
		t.Fatalf("Marshal error with nil fields: %v", err)
	}
}

func TestSession(t *testing.T) {
	now := time.Now()
	expires := now.Add(24 * time.Hour)

	session := Session{
		SessionID: "sess-123",
		Token:     "jwt-token",
		StaffID:   "staff-456",
		CreatedAt: now,
		ExpiresAt: expires,
	}

	if session.SessionID != "sess-123" {
		t.Errorf("SessionID = %q, want %q", session.SessionID, "sess-123")
	}

	if session.ExpiresAt.Before(session.CreatedAt) {
		t.Error("ExpiresAt should be after CreatedAt")
	}
}

func TestErrorResponse(t *testing.T) {
	resp := ErrorResponse{
		Error:   "invalid_credentials",
		Message: "Username or password is incorrect",
		Code:    "AUTH_001",
	}

	if resp.Error != "invalid_credentials" {
		t.Errorf("Error = %q, want %q", resp.Error, "invalid_credentials")
	}

	if resp.Code != "AUTH_001" {
		t.Errorf("Code = %q, want %q", resp.Code, "AUTH_001")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
