package models

import (
	"time"
)

// Session represents a user session (simplified - only session_id and token stored in DB)
type Session struct {
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
}

// SessionCreateRequest represents a session creation request (login)
type SessionCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SessionCreateResponse represents a session creation response
type SessionCreateResponse struct {
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
	Message   string `json:"message"`
	Staff     *Staff `json:"staff,omitempty"`
}

// SessionValidationRequest represents a session validation request
type SessionValidationRequest struct {
	Token string `json:"token"`
}

// SessionValidationResponse represents a session validation response
type SessionValidationResponse struct {
	Valid     bool   `json:"valid"`
	SessionID string `json:"session_id,omitempty"`
	Message   string `json:"message,omitempty"`
	StaffID   string `json:"staff_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Role      string `json:"role,omitempty"`
	FullName  string `json:"full_name,omitempty"`
}

// SessionLogoutRequest represents a session logout request
type SessionLogoutRequest struct {
	Token string `json:"token"`
}

// SessionLogoutResponse represents a session logout response
type SessionLogoutResponse struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id,omitempty"`
	Message   string `json:"message"`
}

// Staff represents a staff member from the database
type Staff struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        *string    `json:"email,omitempty"`
	PasswordHash string     `json:"-"` // Never expose password hash
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// StaffProfile represents staff info for JWT claims
type StaffProfile struct {
	Staff Staff `json:"staff"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
