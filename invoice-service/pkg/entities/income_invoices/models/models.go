package models

import (
	"time"
)

// IncomeInvoice represents an income invoice for customer billing
type IncomeInvoice struct {
	ID               string     `json:"id"`
	OrderID          string     `json:"order_id"`
	PaymentID        *string    `json:"payment_id,omitempty"`
	CustomerID       *string    `json:"customer_id,omitempty"`
	CustomerName     string     `json:"customer_name"`
	CustomerTaxID    *string    `json:"customer_tax_id,omitempty"`
	InvoiceNumber    string     `json:"invoice_number"`
	InvoiceType      string     `json:"invoice_type"`
	Subtotal         float64    `json:"subtotal"`
	TaxAmount        float64    `json:"tax_amount"`
	ServiceCharge    float64    `json:"service_charge"`
	TotalAmount      float64    `json:"total_amount"`
	PaymentMethod    string     `json:"payment_method"`
	XMLData          *string    `json:"xml_data,omitempty"`
	DigitalSignature *string    `json:"digital_signature,omitempty"`
	Status           string     `json:"status"`
	GeneratedAt      *time.Time `json:"generated_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// IncomeInvoiceCreateRequest represents a request to create an income invoice
type IncomeInvoiceCreateRequest struct {
	OrderID          string     `json:"order_id"`
	PaymentID        *string    `json:"payment_id,omitempty"`
	CustomerID       *string    `json:"customer_id,omitempty"`
	CustomerName     string     `json:"customer_name"`
	CustomerTaxID    *string    `json:"customer_tax_id,omitempty"`
	InvoiceNumber    string     `json:"invoice_number"`
	InvoiceType      string     `json:"invoice_type"`
	Subtotal         float64    `json:"subtotal"`
	TaxAmount        float64    `json:"tax_amount"`
	ServiceCharge    float64    `json:"service_charge"`
	TotalAmount      float64    `json:"total_amount"`
	PaymentMethod    string     `json:"payment_method"`
	XMLData          *string    `json:"xml_data,omitempty"`
	DigitalSignature *string    `json:"digital_signature,omitempty"`
	Status           string     `json:"status"`
	GeneratedAt      *time.Time `json:"generated_at,omitempty"`
}

// IncomeInvoiceUpdateRequest represents a request to update an income invoice
type IncomeInvoiceUpdateRequest struct {
	PaymentID        *string    `json:"payment_id,omitempty"`
	CustomerID       *string    `json:"customer_id,omitempty"`
	CustomerName     *string    `json:"customer_name,omitempty"`
	CustomerTaxID    *string    `json:"customer_tax_id,omitempty"`
	InvoiceType      *string    `json:"invoice_type,omitempty"`
	Subtotal         *float64   `json:"subtotal,omitempty"`
	TaxAmount        *float64   `json:"tax_amount,omitempty"`
	ServiceCharge    *float64   `json:"service_charge,omitempty"`
	TotalAmount      *float64   `json:"total_amount,omitempty"`
	PaymentMethod    *string    `json:"payment_method,omitempty"`
	XMLData          *string    `json:"xml_data,omitempty"`
	DigitalSignature *string    `json:"digital_signature,omitempty"`
	Status           *string    `json:"status,omitempty"`
	GeneratedAt      *time.Time `json:"generated_at,omitempty"`
}

// IncomeInvoiceListRequest represents filter parameters for listing income invoices
type IncomeInvoiceListRequest struct {
	CustomerName *string `json:"customer_name,omitempty"`
	InvoiceType  *string `json:"invoice_type,omitempty"`
	Status       *string `json:"status,omitempty"`
	OrderID      *string `json:"order_id,omitempty"`
	CustomerID   *string `json:"customer_id,omitempty"`
	Page         int     `json:"page"`
	Limit        int     `json:"limit"`
}

// IncomeInvoiceListResponse represents a paginated list of income invoices
type IncomeInvoiceListResponse struct {
	Invoices []IncomeInvoice `json:"invoices"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	Limit    int             `json:"limit"`
}
