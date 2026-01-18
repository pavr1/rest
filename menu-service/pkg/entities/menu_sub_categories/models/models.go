package models

import (
	"time"
)

// MenuSubCategory represents a menu sub-category (grouping of menu variants within a category)
type MenuSubCategory struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	ItemType     string    `json:"item_type"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MenuSubCategoryCreateRequest represents a request to create a menu sub-category
type MenuSubCategoryCreateRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	CategoryID   string  `json:"category_id"`
	ItemType     string  `json:"item_type"`
	DisplayOrder int     `json:"display_order"`
	IsActive     bool    `json:"is_active"`
}

// MenuSubCategoryUpdateRequest represents a request to update a menu sub-category
type MenuSubCategoryUpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	CategoryID   *string `json:"category_id,omitempty"`
	ItemType     *string `json:"item_type,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// MenuSubCategoryListRequest represents filter parameters for listing menu sub-categories
type MenuSubCategoryListRequest struct {
	CategoryID *string `json:"category_id,omitempty"`
	ItemType   *string `json:"item_type,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}

// MenuSubCategoryListResponse represents a paginated list of menu sub-categories
type MenuSubCategoryListResponse struct {
	SubCategories []MenuSubCategory `json:"sub_categories"`
	Total         int               `json:"total"`
	Page          int               `json:"page"`
	Limit         int               `json:"limit"`
}
