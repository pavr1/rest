package models

import (
	"time"
)

// StockCount represents an inventory count record for a stock variant
type StockCount struct {
	ID             string    `json:"id"`
	StockVariantID string    `json:"stock_variant_id"`
	InvoiceID      string    `json:"invoice_id"`
	Count          float64   `json:"count"`
	Unit           string    `json:"unit"`
	PurchasedAt    time.Time `json:"purchased_at"`
	IsOut          bool      `json:"is_out"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// Joined fields (optional, populated on list/get)
	StockVariantName *string `json:"stock_variant_name,omitempty"`
	InvoiceNumber    *string `json:"invoice_number,omitempty"`
	SupplierName     *string `json:"supplier_name,omitempty"`
}

// StockCountCreateRequest represents a request to create a stock count record
type StockCountCreateRequest struct {
	StockVariantID string    `json:"stock_variant_id"`
	InvoiceID      string    `json:"invoice_id"`
	Count          float64   `json:"count"`
	Unit           string    `json:"unit"`
	PurchasedAt    time.Time `json:"purchased_at"`
}

// StockCountUpdateRequest represents a request to update a stock count record
type StockCountUpdateRequest struct {
	Count *float64 `json:"count,omitempty"`
	Unit  *string  `json:"unit,omitempty"`
	IsOut *bool    `json:"is_out,omitempty"`
}

// StockCountMarkOutRequest represents a request to mark stock as out
type StockCountMarkOutRequest struct {
	IsOut bool `json:"is_out"`
}

// StockCountListResponse represents a paginated list of stock count records
type StockCountListResponse struct {
	StockCounts []StockCount `json:"stock_counts"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	Limit       int          `json:"limit"`
}
