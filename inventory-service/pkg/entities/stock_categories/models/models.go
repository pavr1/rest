package models

import (
	"time"
)

// StockCategory represents a stock category
type StockCategory struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// StockCategoryCreateRequest represents a request to create a stock category
type StockCategoryCreateRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// StockCategoryUpdateRequest represents a request to update a stock category
type StockCategoryUpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// StockCategoryListResponse represents a paginated list of stock categories
type StockCategoryListResponse struct {
	Categories []StockCategory `json:"categories"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
}
