package models

import (
	"time"
)

// SubMenu represents a sub menu (grouping of menu items within a category)
type SubMenu struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	ImageURL     *string   `json:"image_url,omitempty"`
	ItemType     string    `json:"item_type"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SubMenuCreateRequest represents a request to create a sub menu
type SubMenuCreateRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	CategoryID   string  `json:"category_id"`
	ImageURL     *string `json:"image_url,omitempty"`
	ItemType     string  `json:"item_type"`
	DisplayOrder int     `json:"display_order"`
	IsActive     bool    `json:"is_active"`
}

// SubMenuUpdateRequest represents a request to update a sub menu
type SubMenuUpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	CategoryID   *string `json:"category_id,omitempty"`
	ImageURL     *string `json:"image_url,omitempty"`
	ItemType     *string `json:"item_type,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// SubMenuListRequest represents filter parameters for listing sub menus
type SubMenuListRequest struct {
	CategoryID *string `json:"category_id,omitempty"`
	ItemType   *string `json:"item_type,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}

// SubMenuListResponse represents a paginated list of sub menus
type SubMenuListResponse struct {
	SubMenus []SubMenu `json:"sub_menus"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}
