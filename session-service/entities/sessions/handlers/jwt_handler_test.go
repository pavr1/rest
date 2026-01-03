package handlers

import (
	"session-service/entities/sessions/models"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func newTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	return logger
}

func newTestStaff() *models.Staff {
	return &models.Staff{
		ID:        "test-staff-id-123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Role:      "waiter",
		IsActive:  true,
	}
}

func TestNewJWTHandler(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret", 24*time.Hour, logger)

	if handler == nil {
		t.Fatal("NewJWTHandler returned nil")
	}

	if handler.secretKey != "test-secret" {
		t.Errorf("secretKey = %q, want %q", handler.secretKey, "test-secret")
	}

	if handler.expirationTime != 24*time.Hour {
		t.Errorf("expirationTime = %v, want %v", handler.expirationTime, 24*time.Hour)
	}
}

func TestGenerateSessionID(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret", 24*time.Hour, logger)

	sessionID, err := handler.GenerateSessionID()
	if err != nil {
		t.Fatalf("GenerateSessionID() error = %v", err)
	}

	if len(sessionID) != 32 { // 16 bytes = 32 hex characters
		t.Errorf("sessionID length = %d, want 32", len(sessionID))
	}

	// Test uniqueness
	sessionID2, _ := handler.GenerateSessionID()
	if sessionID == sessionID2 {
		t.Error("GenerateSessionID should generate unique IDs")
	}
}

func TestGenerateToken(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret-key", 1*time.Hour, logger)
	staff := newTestStaff()

	token, expiresAt, err := handler.GenerateToken("session-123", staff)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken returned empty token")
	}

	// Check expiration is approximately 1 hour from now
	expectedExpiry := time.Now().Add(1 * time.Hour)
	diff := expiresAt.Sub(expectedExpiry)
	if diff > time.Second || diff < -time.Second {
		t.Errorf("expiresAt diff from expected = %v, should be < 1 second", diff)
	}
}

func TestValidateToken(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret-key", 1*time.Hour, logger)
	staff := newTestStaff()

	// Generate token
	token, _, err := handler.GenerateToken("session-123", staff)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Validate token
	claims, err := handler.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.StaffID != staff.ID {
		t.Errorf("claims.StaffID = %q, want %q", claims.StaffID, staff.ID)
	}

	if claims.Username != staff.Username {
		t.Errorf("claims.Username = %q, want %q", claims.Username, staff.Username)
	}

	if claims.Role != staff.Role {
		t.Errorf("claims.Role = %q, want %q", claims.Role, staff.Role)
	}

	expectedFullName := "Test User"
	if claims.FullName != expectedFullName {
		t.Errorf("claims.FullName = %q, want %q", claims.FullName, expectedFullName)
	}
}

func TestValidateTokenInvalidSecret(t *testing.T) {
	logger := newTestLogger()
	handler1 := NewJWTHandler("secret-1", 1*time.Hour, logger)
	handler2 := NewJWTHandler("secret-2", 1*time.Hour, logger)
	staff := newTestStaff()

	// Generate token with handler1
	token, _, _ := handler1.GenerateToken("session-123", staff)

	// Try to validate with handler2 (different secret)
	_, err := handler2.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken should fail with wrong secret")
	}
}

func TestValidateTokenInvalidFormat(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret", 1*time.Hour, logger)

	_, err := handler.ValidateToken("invalid-token-format")
	if err == nil {
		t.Error("ValidateToken should fail for invalid token format")
	}
}

func TestValidateTokenEmpty(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret", 1*time.Hour, logger)

	_, err := handler.ValidateToken("")
	if err == nil {
		t.Error("ValidateToken should fail for empty token")
	}
}

func TestGenerateTokenHash(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret", 1*time.Hour, logger)

	token := "some-jwt-token-string"
	hash := handler.GenerateTokenHash(token)

	if len(hash) != 64 { // SHA256 = 32 bytes = 64 hex chars
		t.Errorf("hash length = %d, want 64", len(hash))
	}

	// Same input should produce same hash
	hash2 := handler.GenerateTokenHash(token)
	if hash != hash2 {
		t.Error("Same token should produce same hash")
	}

	// Different input should produce different hash
	hash3 := handler.GenerateTokenHash("different-token")
	if hash == hash3 {
		t.Error("Different tokens should produce different hashes")
	}
}

func TestGetTokenExpiration(t *testing.T) {
	logger := newTestLogger()
	handler := NewJWTHandler("test-secret", 2*time.Hour, logger)
	staff := newTestStaff()

	token, expectedExpiry, _ := handler.GenerateToken("session-123", staff)

	expiry, err := handler.GetTokenExpiration(token)
	if err != nil {
		t.Fatalf("GetTokenExpiration() error = %v", err)
	}

	diff := expiry.Sub(expectedExpiry)
	if diff > time.Second || diff < -time.Second {
		t.Errorf("expiry diff = %v, should be < 1 second", diff)
	}
}

func TestGetExpirationTime(t *testing.T) {
	logger := newTestLogger()
	expectedDuration := 12 * time.Hour
	handler := NewJWTHandler("test-secret", expectedDuration, logger)

	if handler.GetExpirationTime() != expectedDuration {
		t.Errorf("GetExpirationTime() = %v, want %v", handler.GetExpirationTime(), expectedDuration)
	}
}
