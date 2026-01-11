package models

import (
	"time"
)

// InvoiceDetail represents a line item in a purchase invoice
type InvoiceDetail struct {
	ID              string   `json:"id"`
	InvoiceID       string   `json:"invoice_id"`
	StockItemID     *string  `json:"stock_item_id,omitempty"`
	StockItemName   *string  `json:"stock_item_name,omitempty"`
	Description     string   `json:"description"`
	Quantity        float64  `json:"quantity"`
	UnitOfMeasure   string   `json:"unit_of_measure"`
	ItemsPerUnit    *float64 `json:"items_per_unit,omitempty"`
	UnitPrice       float64  `json:"unit_price"`
	TotalPrice      float64  `json:"total_price"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	BatchNumber     *string  `json:"batch_number,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// InvoiceDetailCreateRequest represents a request to create an invoice detail
type InvoiceDetailCreateRequest struct {
	StockItemID   *string   `json:"stock_item_id,omitempty"`
	Description   string    `json:"description"`
	Quantity      float64   `json:"quantity"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	ItemsPerUnit  *float64  `json:"items_per_unit,omitempty"`
	UnitPrice     float64   `json:"unit_price"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	BatchNumber   *string   `json:"batch_number,omitempty"`
}

// InvoiceDetailUpdateRequest represents a request to update an invoice detail
type InvoiceDetailUpdateRequest struct {
	StockItemID   *string   `json:"stock_item_id,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Quantity      *float64  `json:"quantity,omitempty"`
	UnitOfMeasure *string   `json:"unit_of_measure,omitempty"`
	ItemsPerUnit  *float64  `json:"items_per_unit,omitempty"`
	UnitPrice     *float64  `json:"unit_price,omitempty"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	BatchNumber   *string   `json:"batch_number,omitempty"`
}

// InvoiceDetailListResponse represents a list of invoice details for an invoice
type InvoiceDetailListResponse struct {
	InvoiceID      string          `json:"invoice_id"`
	InvoiceDetails []InvoiceDetail `json:"invoice_details"`
}