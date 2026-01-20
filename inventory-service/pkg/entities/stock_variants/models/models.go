package models

import (
	"time"
)

// StockVariant represents a stock variant
type StockVariant struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	StockSubCategoryID string    `json:"stock_sub_category_id"`
	InvoiceID          *string   `json:"invoice_id,omitempty"`
	Unit               string    `json:"unit"`
	NumberOfUnits      float64   `json:"number_of_units"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// StockVariantCreateRequest represents a request to create a stock variant
type StockVariantCreateRequest struct {
	Name               string  `json:"name"`
	StockSubCategoryID string  `json:"stock_sub_category_id"`
	InvoiceID          *string `json:"invoice_id,omitempty"`
	Unit               string  `json:"unit"`
	NumberOfUnits      float64 `json:"number_of_units"`
	IsActive           *bool   `json:"is_active,omitempty"`
}

// StockVariantUpdateRequest represents a request to update a stock variant
type StockVariantUpdateRequest struct {
	Name          *string  `json:"name,omitempty"`
	Unit          *string  `json:"unit,omitempty"`
	NumberOfUnits *float64 `json:"number_of_units,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

// StockVariantListResponse represents a paginated list of stock variants
type StockVariantListResponse struct {
	Variants []StockVariant `json:"variants"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
}
