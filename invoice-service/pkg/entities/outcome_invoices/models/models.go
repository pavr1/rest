package models

import (
	"time"

	invoiceItemModels "invoice-service/pkg/entities/invoice_items/models"
)

// InvoiceItem alias for easier reference
type InvoiceItem = invoiceItemModels.InvoiceItem
type InvoiceItemCreateRequest = invoiceItemModels.InvoiceItemCreateRequest

// OutcomeInvoice represents an outcome invoice for supplier purchases
type OutcomeInvoice struct {
	ID             string        `json:"id"`
	InvoiceNumber  string        `json:"invoice_number"`
	SupplierID     *string       `json:"supplier_id,omitempty"`
	TransactionDate time.Time     `json:"transaction_date"`
	TotalAmount    *float64      `json:"total_amount,omitempty"`
	ImageURL       *string       `json:"image_url,omitempty"`
	Notes          *string       `json:"notes,omitempty"`
	InvoiceItems   []InvoiceItem `json:"invoice_items,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// OutcomeInvoiceCreateRequest represents a request to create an outcome invoice
type OutcomeInvoiceCreateRequest struct {
	InvoiceNumber  string                     `json:"invoice_number"`
	SupplierID     *string                    `json:"supplier_id,omitempty"`
	TransactionDate time.Time                  `json:"transaction_date"`
	TotalAmount    *float64                   `json:"total_amount,omitempty"`
	ImageURL       *string                    `json:"image_url,omitempty"`
	Notes          *string                    `json:"notes,omitempty"`
	InvoiceItems   []InvoiceItemCreateRequest `json:"invoice_items,omitempty"`
}

// OutcomeInvoiceUpdateRequest represents a request to update an outcome invoice
type OutcomeInvoiceUpdateRequest struct {
	SupplierID     *string    `json:"supplier_id,omitempty"`
	TransactionDate *time.Time `json:"transaction_date,omitempty"`
	TotalAmount    *float64   `json:"total_amount,omitempty"`
	ImageURL       *string    `json:"image_url,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
}

// OutcomeInvoiceListRequest represents filter parameters for listing outcome invoices
type OutcomeInvoiceListRequest struct {
	SupplierID *string `json:"supplier_id,omitempty"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}

// OutcomeInvoiceListResponse represents a paginated list of outcome invoices
type OutcomeInvoiceListResponse struct {
	Invoices []OutcomeInvoice `json:"invoices"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}