package models

import (
	"time"
)

// Supplier represents a supplier in the system
type Supplier struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ContactName *string   `json:"contact_name,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Address     *string   `json:"address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SupplierCreateRequest represents a request to create a supplier
type SupplierCreateRequest struct {
	Name        string  `json:"name"`
	ContactName *string `json:"contact_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Email       *string `json:"email,omitempty"`
	Address     *string `json:"address,omitempty"`
}

// SupplierUpdateRequest represents a request to update a supplier
type SupplierUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	ContactName *string `json:"contact_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Email       *string `json:"email,omitempty"`
	Address     *string `json:"address,omitempty"`
}

// SupplierListRequest represents filter parameters for listing suppliers
type SupplierListRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}

// SupplierListResponse represents a paginated list of suppliers
type SupplierListResponse struct {
	Suppliers []Supplier `json:"suppliers"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
