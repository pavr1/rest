package models

import (
	"fmt"
	"strings"
	"time"
)

// StockCount represents an inventory count record for a stock variant
type StockCount struct {
	ID              string    `json:"id"`
	StockVariantID  string    `json:"stock_variant_id"`
	InvoiceID       *string   `json:"invoice_id,omitempty"`
	Count           float64   `json:"count"`
	Unit            string    `json:"unit"`
	UnitPrice       *float64  `json:"unit_price,omitempty"`
	CostPerPortion  *float64  `json:"cost_per_portion,omitempty"`
	PurchasedAt     time.Time `json:"purchased_at"`
	IsOut           bool      `json:"is_out"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Joined fields (optional, populated on list/get)
	StockVariantName *string `json:"stock_variant_name,omitempty"`
	InvoiceNumber    *string `json:"invoice_number,omitempty"`
	SupplierName     *string `json:"supplier_name,omitempty"`
}

// StockCountCreateRequest represents a request to create a stock count record
type StockCountCreateRequest struct {
	StockVariantID string    `json:"stock_variant_id"`
	InvoiceID      *string   `json:"invoice_id,omitempty"`
	Count          float64   `json:"count"`
	Unit           string    `json:"unit"`
	UnitPrice      *float64  `json:"unit_price,omitempty"`
	PurchasedAt    time.Time `json:"purchased_at"`
}

// StockCountUpdateRequest represents a request to update a stock count record
type StockCountUpdateRequest struct {
	Count     *float64 `json:"count,omitempty"`
	Unit      *string  `json:"unit,omitempty"`
	UnitPrice *float64 `json:"unit_price,omitempty"`
	IsOut     *bool    `json:"is_out,omitempty"`
}

// Supported units for stock count (all convertible to kg)
const (
	UnitKG = "kg"  // Kilograms
	UnitG  = "g"   // Grams
	UnitL  = "l"   // Liters (treated as kg for cost calculation)
	UnitML = "ml"  // Milliliters (treated as g for cost calculation)
)

// ConvertToKG converts the given count and unit to kilograms
func ConvertToKG(count float64, unit string) (float64, error) {
	switch strings.ToLower(unit) {
	case UnitKG:
		return count, nil
	case UnitG:
		return count / 1000, nil
	case UnitL:
		// Treat 1 liter as approximately 1 kg
		return count, nil
	case UnitML:
		// Treat 1 ml as approximately 1 g
		return count / 1000, nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s (supported: kg, g, l, ml)", unit)
	}
}

// CalculateCostPerPortion calculates the cost per portion
// totalKG: total weight in kilograms
// unitPrice: price per unit of the original count
// count: original count (before conversion to kg)
// portionGrams: default portion size in grams (from settings)
func CalculateCostPerPortion(totalKG float64, unitPrice float64, portionGrams float64) float64 {
	if totalKG <= 0 || portionGrams <= 0 {
		return 0
	}
	// Total price = unitPrice (this is the total price for the purchase, not per kg)
	// Number of portions = totalKG / (portionGrams / 1000)
	portionKG := portionGrams / 1000
	numPortions := totalKG / portionKG
	if numPortions <= 0 {
		return 0
	}
	return unitPrice / numPortions
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
