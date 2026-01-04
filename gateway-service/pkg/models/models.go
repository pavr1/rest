package models

// SessionValidationRequest represents a session validation request
type SessionValidationRequest struct {
	SessionID string `json:"session_id"`
}

// SessionValidationResponse represents a session validation response
type SessionValidationResponse struct {
	Valid       bool     `json:"valid"`
	SessionID   string   `json:"session_id,omitempty"`
	Message     string   `json:"message,omitempty"`
	UserID      string   `json:"user_id,omitempty"`
	Username    string   `json:"username,omitempty"`
	RoleName    string   `json:"role_name,omitempty"`
	FullName    string   `json:"full_name,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// SessionCreateRequest represents a session creation request
type SessionCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SessionCreateResponse represents a session creation response
type SessionCreateResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// SessionLogoutRequest represents a session logout request
type SessionLogoutRequest struct {
	SessionID string `json:"session_id"`
}
