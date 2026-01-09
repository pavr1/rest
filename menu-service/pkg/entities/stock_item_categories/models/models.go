package models

import (
	"time"
)

// StockItemCategory represents a stock item category
type StockItemCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockItemCategoryCreateRequest represents a request to create a stock item category
type StockItemCategoryCreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// StockItemCategoryUpdateRequest represents a request to update a stock item category
type StockItemCategoryUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// StockItemCategoryListResponse represents a paginated list of stock item categories
type StockItemCategoryListResponse struct {
	Categories []StockItemCategory `json:"categories"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
}
