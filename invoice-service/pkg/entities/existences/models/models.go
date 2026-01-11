package models

import (
	"time"
)

// Existence represents a stock batch created from an invoice detail
type Existence struct {
	ID                 string    `json:"id"`
	InvoiceDetailID    string    `json:"invoice_detail_id"`
	StockItemID        string    `json:"stock_item_id"`
	StockItemName      string    `json:"stock_item_name,omitempty"`
	UnitsPurchased     float64   `json:"units_purchased"`
	CostPerUnit        float64   `json:"cost_per_unit"`
	TotalCost          float64   `json:"total_cost"`
	ExpiryDate         *time.Time `json:"expiry_date,omitempty"`
	BatchNumber        *string   `json:"batch_number,omitempty"`
	CurrentStock       float64   `json:"current_stock"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ExistenceCreateRequest represents a request to create an existence from an invoice detail
type ExistenceCreateRequest struct {
	InvoiceDetailID string     `json:"invoice_detail_id"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	BatchNumber     *string    `json:"batch_number,omitempty"`
}

// ExistenceUpdateRequest represents a request to update an existence
type ExistenceUpdateRequest struct {
	CurrentStock *float64   `json:"current_stock,omitempty"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	BatchNumber  *string    `json:"batch_number,omitempty"`
}

// ExistenceListRequest represents filter parameters for listing existences
type ExistenceListRequest struct {
	StockItemID *string `json:"stock_item_id,omitempty"`
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
}

// ExistenceListResponse represents a paginated list of existences
type ExistenceListResponse struct {
	Existences []Existence `json:"existences"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}