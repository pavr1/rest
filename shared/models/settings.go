package models

import "time"

// Setting represents a configuration setting stored in the database
type Setting struct {
	SettingID   string    `json:"setting_id"`
	Service     string    `json:"service"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetSettingsByServiceRequest is the request body for getting settings by service
type GetSettingsByServiceRequest struct {
	Service string `json:"service"`
}

// GetSettingsByKeyRequest is the request body for getting a setting by key
type GetSettingsByKeyRequest struct {
	Service string `json:"service"`
	Key     string `json:"key"`
}

// UpdateSettingRequest is the request body for updating a setting
type UpdateSettingRequest struct {
	Service string `json:"service"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

// SettingsResponse is the response structure for settings endpoints
type SettingsResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

