package models

import (
	"time"
)

// PurchaseInvoice represents a purchase invoice from a supplier
type PurchaseInvoice struct {
	ID            string     `json:"id"`
	InvoiceNumber string     `json:"invoice_number"`
	SupplierName  string     `json:"supplier_name"`
	InvoiceDate   time.Time  `json:"invoice_date"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	TotalAmount   *float64   `json:"total_amount,omitempty"`
	Status        string     `json:"status"`
	ImageURL      *string    `json:"image_url,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// PurchaseInvoiceCreateRequest represents a request to create a purchase invoice
type PurchaseInvoiceCreateRequest struct {
	InvoiceNumber string     `json:"invoice_number"`
	SupplierName  string     `json:"supplier_name"`
	InvoiceDate   time.Time  `json:"invoice_date"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	TotalAmount   *float64   `json:"total_amount,omitempty"`
	Status        string     `json:"status"`
	ImageURL      *string    `json:"image_url,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
}

// PurchaseInvoiceUpdateRequest represents a request to update a purchase invoice
type PurchaseInvoiceUpdateRequest struct {
	SupplierName *string    `json:"supplier_name,omitempty"`
	InvoiceDate  *time.Time `json:"invoice_date,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	TotalAmount  *float64   `json:"total_amount,omitempty"`
	Status       *string    `json:"status,omitempty"`
	ImageURL     *string    `json:"image_url,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
}

// PurchaseInvoiceListRequest represents filter parameters for listing purchase invoices
type PurchaseInvoiceListRequest struct {
	SupplierName *string `json:"supplier_name,omitempty"`
	Status       *string `json:"status,omitempty"`
	Page         int     `json:"page"`
	Limit        int     `json:"limit"`
}

// PurchaseInvoiceListResponse represents a paginated list of purchase invoices
type PurchaseInvoiceListResponse struct {
	Invoices []PurchaseInvoice `json:"invoices"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}