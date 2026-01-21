package models

import (
	"time"
)

// StockVariant represents a stock variant (defines the item type, actual counts are in stock_count)
type StockVariant struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        *string   `json:"description,omitempty"`
	StockSubCategoryID string    `json:"stock_sub_category_id"`
	AvgCost            float64   `json:"avg_cost"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// StockVariantCreateRequest represents a request to create a stock variant
type StockVariantCreateRequest struct {
	Name               string  `json:"name"`
	Description        *string `json:"description,omitempty"`
	StockSubCategoryID string  `json:"stock_sub_category_id"`
	IsActive           *bool   `json:"is_active,omitempty"`
}

// StockVariantUpdateRequest represents a request to update a stock variant
type StockVariantUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// StockVariantListResponse represents a paginated list of stock variants
type StockVariantListResponse struct {
	Variants []StockVariant `json:"variants"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
}
