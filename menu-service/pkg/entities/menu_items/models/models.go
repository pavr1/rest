package models

import (
	"encoding/json"
	"time"
)

// MenuItem represents a menu item
type MenuItem struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    *string         `json:"description,omitempty"`
	CategoryID     string          `json:"category_id"`
	CategoryName   string          `json:"category_name,omitempty"`
	Price          float64         `json:"price"`
	ItemCost       *float64        `json:"item_cost,omitempty"`
	HappyHourPrice *float64        `json:"happy_hour_price,omitempty"`
	ImageURL       *string         `json:"image_url,omitempty"`
	IsAvailable    bool            `json:"is_available"`
	ItemType       string          `json:"item_type"`
	MenuTypes      json.RawMessage `json:"menu_types"`
	DietaryTags    json.RawMessage `json:"dietary_tags,omitempty"`
	Allergens      json.RawMessage `json:"allergens,omitempty"`
	IsAlcoholic    bool            `json:"is_alcoholic"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// MenuItemCreateRequest represents a request to create a menu item
type MenuItemCreateRequest struct {
	Name           string          `json:"name"`
	Description    *string         `json:"description,omitempty"`
	CategoryID     string          `json:"category_id"`
	Price          float64         `json:"price"`
	HappyHourPrice *float64        `json:"happy_hour_price,omitempty"`
	ImageURL       *string         `json:"image_url,omitempty"`
	IsAvailable    bool            `json:"is_available"`
	ItemType       string          `json:"item_type"`
	MenuTypes      json.RawMessage `json:"menu_types"`
	DietaryTags    json.RawMessage `json:"dietary_tags,omitempty"`
	Allergens      json.RawMessage `json:"allergens,omitempty"`
	IsAlcoholic    bool            `json:"is_alcoholic"`
}

// MenuItemUpdateRequest represents a request to update a menu item
type MenuItemUpdateRequest struct {
	Name           *string          `json:"name,omitempty"`
	Description    *string          `json:"description,omitempty"`
	CategoryID     *string          `json:"category_id,omitempty"`
	Price          *float64         `json:"price,omitempty"`
	HappyHourPrice *float64         `json:"happy_hour_price,omitempty"`
	ImageURL       *string          `json:"image_url,omitempty"`
	IsAvailable    *bool            `json:"is_available,omitempty"`
	ItemType       *string          `json:"item_type,omitempty"`
	MenuTypes      *json.RawMessage `json:"menu_types,omitempty"`
	DietaryTags    *json.RawMessage `json:"dietary_tags,omitempty"`
	Allergens      *json.RawMessage `json:"allergens,omitempty"`
	IsAlcoholic    *bool            `json:"is_alcoholic,omitempty"`
}

// MenuItemAvailabilityRequest represents a request to toggle availability
type MenuItemAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

// MenuItemListRequest represents filter parameters for listing menu items
type MenuItemListRequest struct {
	MenuType    *string `json:"menu_type,omitempty"`
	CategoryID  *string `json:"category_id,omitempty"`
	IsAvailable *bool   `json:"is_available,omitempty"`
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
}

// MenuItemListResponse represents a paginated list of menu items
type MenuItemListResponse struct {
	Items []MenuItem `json:"items"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}
