package models

import (
	"time"
)

// StockSubCategory represents a stock sub-category
type StockSubCategory struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	StockCategoryID string  `json:"stock_category_id"`
	DisplayOrder  int       `json:"display_order"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// StockSubCategoryCreateRequest represents a request to create a stock sub-category
type StockSubCategoryCreateRequest struct {
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	StockCategoryID string  `json:"stock_category_id"`
	DisplayOrder    *int    `json:"display_order,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
}

// StockSubCategoryUpdateRequest represents a request to update a stock sub-category
type StockSubCategoryUpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// StockSubCategoryListResponse represents a paginated list of stock sub-categories
type StockSubCategoryListResponse struct {
	SubCategories []StockSubCategory `json:"sub_categories"`
	Total         int                 `json:"total"`
	Page          int                 `json:"page"`
	Limit         int                 `json:"limit"`
}
