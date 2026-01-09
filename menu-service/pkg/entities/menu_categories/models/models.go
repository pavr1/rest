package models

import (
	"time"
)

// MenuCategory represents a menu category
type MenuCategory struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"display_order"`
	Description  *string   `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MenuCategoryCreateRequest represents a request to create a menu category
type MenuCategoryCreateRequest struct {
	Name         string  `json:"name"`
	DisplayOrder int     `json:"display_order"`
	Description  *string `json:"description,omitempty"`
}

// MenuCategoryUpdateRequest represents a request to update a menu category
type MenuCategoryUpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	Description  *string `json:"description,omitempty"`
}

// MenuCategoryListResponse represents a paginated list of menu categories
type MenuCategoryListResponse struct {
	Categories []MenuCategory `json:"categories"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}
