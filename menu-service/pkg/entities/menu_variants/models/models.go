package models

import (
	"encoding/json"
	"time"
)

// MenuVariant represents a menu variant (actual orderable item with pricing)
type MenuVariant struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Description       *string         `json:"description,omitempty"`
	SubCategoryID     string          `json:"sub_category_id"`
	SubCategoryName   string          `json:"sub_category_name,omitempty"`
	ItemType          string          `json:"item_type,omitempty"` // Inherited from sub_category
	Price             float64         `json:"price"`
	ItemCost          *float64        `json:"item_cost,omitempty"`
	HappyHourPrice    *float64        `json:"happy_hour_price,omitempty"`
	ImageURL          *string         `json:"image_url,omitempty"`
	IsAvailable       bool            `json:"is_available"`
	PreparationTime   *int            `json:"preparation_time,omitempty"`
	MenuTypes         json.RawMessage `json:"menu_types"`
	DietaryTags       json.RawMessage `json:"dietary_tags,omitempty"`
	Allergens         json.RawMessage `json:"allergens,omitempty"`
	IsAlcoholic       bool            `json:"is_alcoholic"`
	DisplayOrder      int             `json:"display_order"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// MenuVariantCreateRequest represents a request to create a menu item
type MenuVariantCreateRequest struct {
	Name            string          `json:"name"`
	Description     *string         `json:"description,omitempty"`
	SubCategoryID       string          `json:"sub_category_id"`
	Price           float64         `json:"price"`
	HappyHourPrice  *float64        `json:"happy_hour_price,omitempty"`
	ImageURL        *string         `json:"image_url,omitempty"`
	IsAvailable     bool            `json:"is_available"`
	PreparationTime *int            `json:"preparation_time,omitempty"`
	MenuTypes       json.RawMessage `json:"menu_types"`
	DietaryTags     json.RawMessage `json:"dietary_tags,omitempty"`
	Allergens       json.RawMessage `json:"allergens,omitempty"`
	IsAlcoholic     bool            `json:"is_alcoholic"`
	DisplayOrder    int             `json:"display_order"`
}

// MenuVariantUpdateRequest represents a request to update a menu item
type MenuVariantUpdateRequest struct {
	Name            *string          `json:"name,omitempty"`
	Description     *string          `json:"description,omitempty"`
	SubCategoryID       *string          `json:"sub_category_id,omitempty"`
	Price           *float64         `json:"price,omitempty"`
	HappyHourPrice  *float64         `json:"happy_hour_price,omitempty"`
	ImageURL        *string          `json:"image_url,omitempty"`
	IsAvailable     *bool            `json:"is_available,omitempty"`
	PreparationTime *int             `json:"preparation_time,omitempty"`
	MenuTypes       *json.RawMessage `json:"menu_types,omitempty"`
	DietaryTags     *json.RawMessage `json:"dietary_tags,omitempty"`
	Allergens       *json.RawMessage `json:"allergens,omitempty"`
	IsAlcoholic     *bool            `json:"is_alcoholic,omitempty"`
	DisplayOrder    *int             `json:"display_order,omitempty"`
}

// MenuVariantAvailabilityRequest represents a request to toggle availability
type MenuVariantAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

// MenuVariantListRequest represents filter parameters for listing menu items
type MenuVariantListRequest struct {
	MenuType    *string `json:"menu_type,omitempty"`
	SubCategoryID   *string `json:"sub_category_id,omitempty"`
	IsAvailable *bool   `json:"is_available,omitempty"`
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
}

// MenuVariantListResponse represents a paginated list of menu items
type MenuVariantListResponse struct {
	Items []MenuVariant `json:"items"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}
