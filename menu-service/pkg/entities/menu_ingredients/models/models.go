package models

import (
	"time"
)

// MenuIngredient represents a menu variant's stock variant requirement
type MenuIngredient struct {
	ID               string    `json:"id"`
	MenuVariantID    string    `json:"menu_variant_id"`
	StockVariantID   string    `json:"stock_variant_id"`
	StockVariantName string    `json:"stock_variant_name,omitempty"`
	Quantity         float64   `json:"quantity"`
	IsOptional       bool      `json:"is_optional"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// MenuIngredientCreateRequest represents a request to add an ingredient to a menu variant
type MenuIngredientCreateRequest struct {
	StockVariantID string  `json:"stock_variant_id"`
	Quantity       float64 `json:"quantity"`
	IsOptional     bool    `json:"is_optional,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}

// MenuIngredientUpdateRequest represents a request to update an ingredient
type MenuIngredientUpdateRequest struct {
	Quantity   *float64 `json:"quantity,omitempty"`
	IsOptional *bool    `json:"is_optional,omitempty"`
	Notes      *string  `json:"notes,omitempty"`
}

// MenuIngredientListResponse represents a list of ingredients for a menu variant
type MenuIngredientListResponse struct {
	MenuVariantID string           `json:"menu_variant_id"`
	Ingredients   []MenuIngredient `json:"ingredients"`
}
