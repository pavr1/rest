package models

import (
	"time"
)

// ValidUnits contains the predefined list of valid units
var ValidUnits = []string{"kg", "g", "lb", "oz", "l", "ml", "unit", "dozen"}

// StockItem represents a stock item
type StockItem struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Unit         string   `json:"unit"`
	Description  *string  `json:"description,omitempty"`
	CategoryID   *string  `json:"category_id,omitempty"`
	CategoryName *string  `json:"category_name,omitempty"`
	CurrentStock float64  `json:"current_stock"`
	UnitCost     *float64 `json:"unit_cost,omitempty"`
	TotalValue   float64  `json:"total_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// StockItemCreateRequest represents a request to create a stock item
type StockItemCreateRequest struct {
	Name        string   `json:"name"`
	Unit        string   `json:"unit"`
	Description *string  `json:"description,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
	UnitCost    *float64 `json:"unit_cost,omitempty"`
}

// StockItemUpdateRequest represents a request to update a stock item
type StockItemUpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Unit        *string  `json:"unit,omitempty"`
	Description *string  `json:"description,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
	UnitCost    *float64 `json:"unit_cost,omitempty"`
}

// StockItemListRequest represents filter parameters for listing stock items
type StockItemListRequest struct {
	CategoryID *string `json:"category_id,omitempty"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}

// StockItemListResponse represents a paginated list of stock items
type StockItemListResponse struct {
	Items []StockItem `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

// IsValidUnit checks if a unit is in the predefined list
func IsValidUnit(unit string) bool {
	for _, validUnit := range ValidUnits {
		if validUnit == unit {
			return true
		}
	}
	return false
}
