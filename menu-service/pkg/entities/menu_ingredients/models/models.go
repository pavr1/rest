package models

import (
	"fmt"
	"time"
)

// MenuIngredient represents a menu variant's ingredient (either stock variant or menu sub-category)
type MenuIngredient struct {
	ID                    string    `json:"id"`
	MenuVariantID         string    `json:"menu_variant_id"`
	StockVariantID        *string   `json:"stock_variant_id,omitempty"`
	StockVariantName      string    `json:"stock_variant_name,omitempty"`
	MenuSubCategoryID     *string   `json:"menu_sub_category_id,omitempty"`
	MenuSubCategoryName   string    `json:"menu_sub_category_name,omitempty"`
	Quantity              float64   `json:"quantity"`
	IsOptional            bool      `json:"is_optional"`
	Notes                 *string   `json:"notes,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// IngredientType returns "stock" or "menu" based on which reference is set
func (m *MenuIngredient) IngredientType() string {
	if m.StockVariantID != nil {
		return "stock"
	}
	return "menu"
}

// MenuIngredientCreateRequest represents a request to add an ingredient to a menu variant
type MenuIngredientCreateRequest struct {
	StockVariantID    *string `json:"stock_variant_id,omitempty"`
	MenuSubCategoryID *string `json:"menu_sub_category_id,omitempty"`
	Quantity          float64 `json:"quantity"`
	IsOptional        bool    `json:"is_optional,omitempty"`
	Notes             *string `json:"notes,omitempty"`
}

// Validate ensures exactly one of StockVariantID or MenuSubCategoryID is provided
func (r *MenuIngredientCreateRequest) Validate() error {
	hasStock := r.StockVariantID != nil && *r.StockVariantID != ""
	hasMenu := r.MenuSubCategoryID != nil && *r.MenuSubCategoryID != ""

	if hasStock && hasMenu {
		return fmt.Errorf("cannot specify both stock_variant_id and menu_sub_category_id")
	}
	if !hasStock && !hasMenu {
		return fmt.Errorf("must specify either stock_variant_id or menu_sub_category_id")
	}
	return nil
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
