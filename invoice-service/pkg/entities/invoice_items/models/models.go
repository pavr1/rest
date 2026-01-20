package models

import (
	"time"
)

// InvoiceItem represents a line item in an invoice (can be for outcome or income invoices)
type InvoiceItem struct {
	ID             string     `json:"id"`
	InvoiceID      string     `json:"invoice_id"`
	StockVariantID *string    `json:"stock_variant_id,omitempty"`
	Detail         string     `json:"detail,omitempty"`
	Count          float64    `json:"count"`
	UnitType       string     `json:"unit_type"`
	Price          float64    `json:"price"`
	ItemsPerUnit   int        `json:"items_per_unit"`
	Total          float64    `json:"total"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// Joined fields (optional, populated on list/get)
	StockVariantName *string `json:"stock_variant_name,omitempty"`
	// Deprecated fields - kept for backward compatibility
	InventoryCategoryID    *string `json:"inventory_category_id,omitempty"`
	InventorySubCategoryID *string `json:"inventory_sub_category_id,omitempty"`
}

// InvoiceItemCreateRequest represents a request to create an invoice item
type InvoiceItemCreateRequest struct {
	InvoiceID      string     `json:"invoice_id"`
	StockVariantID *string    `json:"stock_variant_id,omitempty"`
	Detail         string     `json:"detail,omitempty"`
	Count          float64    `json:"count"`
	UnitType       string     `json:"unit_type"`
	Price          float64    `json:"price"`
	ItemsPerUnit   int        `json:"items_per_unit"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}

// InvoiceItemUpdateRequest represents a request to update an invoice item
type InvoiceItemUpdateRequest struct {
	StockVariantID *string    `json:"stock_variant_id,omitempty"`
	Detail         *string    `json:"detail,omitempty"`
	Count          *float64   `json:"count,omitempty"`
	UnitType       *string    `json:"unit_type,omitempty"`
	Price          *float64   `json:"price,omitempty"`
	ItemsPerUnit   *int       `json:"items_per_unit,omitempty"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}

// InvoiceItemListRequest represents filter parameters for listing invoice items
type InvoiceItemListRequest struct {
	InvoiceID *string `json:"invoice_id,omitempty"`
	Page      int     `json:"page"`
	Limit     int     `json:"limit"`
}

// InvoiceItemListResponse represents a paginated list of invoice items
type InvoiceItemListResponse struct {
	Items []InvoiceItem `json:"items"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}
