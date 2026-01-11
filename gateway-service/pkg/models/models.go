package models

// TokenValidationRequest represents a token validation request
type TokenValidationRequest struct {
	Token string `json:"token"`
}

// TokenValidationResponse represents a token validation response
type TokenValidationResponse struct {
	Valid       bool     `json:"valid"`
	Token       string   `json:"token,omitempty"`
	Message     string   `json:"message,omitempty"`
	StaffID     string   `json:"staff_id,omitempty"`
	Username    string   `json:"username,omitempty"`
	Role        string   `json:"role,omitempty"`
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
	Token string `json:"token"`
}
