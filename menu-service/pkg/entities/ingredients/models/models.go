package models

import (
	"time"
)

// Ingredient represents a menu item's stock item relationship
type Ingredient struct {
	ID            string    `json:"id"`
	MenuItemID    string    `json:"menu_item_id"`
	StockItemID   string    `json:"stock_item_id"`
	StockItemName string    `json:"stock_item_name,omitempty"`
	StockItemUnit string    `json:"stock_item_unit,omitempty"`
	Quantity      float64   `json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// IngredientCreateRequest represents a request to add an ingredient to a menu item
type IngredientCreateRequest struct {
	StockItemID string  `json:"stock_item_id"`
	Quantity    float64 `json:"quantity"`
}

// IngredientUpdateRequest represents a request to update an ingredient quantity
type IngredientUpdateRequest struct {
	Quantity float64 `json:"quantity"`
}

// IngredientListResponse represents a list of ingredients for a menu item
type IngredientListResponse struct {
	MenuItemID  string       `json:"menu_item_id"`
	Ingredients []Ingredient `json:"ingredients"`
}

// MenuItemCost represents the calculated cost for a menu item
type MenuItemCost struct {
	MenuItemID     string            `json:"menu_item_id"`
	MenuItemName   string            `json:"menu_item_name"`
	TotalCost      float64           `json:"total_cost"`
	IngredientCost []IngredientCost  `json:"ingredient_costs"`
	CalculatedAt   time.Time         `json:"calculated_at"`
}

// IngredientCost represents the cost breakdown for one ingredient
type IngredientCost struct {
	StockItemID   string  `json:"stock_item_id"`
	StockItemName string  `json:"stock_item_name"`
	Quantity      float64 `json:"quantity"`
	UnitCost      float64 `json:"unit_cost"`
	TotalCost     float64 `json:"total_cost"`
}
