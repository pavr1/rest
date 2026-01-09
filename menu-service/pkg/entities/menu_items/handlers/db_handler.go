package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"menu-service/pkg/entities/menu_items/models"
	menuItemSQL "menu-service/pkg/entities/menu_items/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for menu items
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *menuItemSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := menuItemSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of menu items
func (h *DBHandler) List(req *models.MenuItemListRequest) (*models.MenuItemListResponse, error) {
	offset := (req.Page - 1) * req.Limit

	// Prepare menu_types filter as JSONB
	var menuTypesJSON interface{}
	if req.MenuType != nil {
		menuTypesJSON = fmt.Sprintf(`["%s"]`, *req.MenuType)
	}

	countQuery, err := h.queries.Get(menuItemSQL.CountMenuItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, req.CategoryID, req.IsAvailable, menuTypesJSON).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count menu items: %w", err)
	}

	listQuery, err := h.queries.Get(menuItemSQL.ListMenuItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, req.CategoryID, req.IsAvailable, menuTypesJSON, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items: %w", err)
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		item, err := h.scanMenuItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return &models.MenuItemListResponse{
		Items: items,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// GetByID returns a menu item by ID
func (h *DBHandler) GetByID(id string) (*models.MenuItem, error) {
	query, err := h.queries.Get(menuItemSQL.GetMenuItemByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id)
	return h.scanMenuItemRow(row)
}

// Create creates a new menu item
func (h *DBHandler) Create(req *models.MenuItemCreateRequest) (*models.MenuItem, error) {
	query, err := h.queries.Get(menuItemSQL.CreateMenuItemQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query,
		req.Name, req.Description, req.CategoryID, req.Price, req.HappyHourPrice,
		req.ImageURL, req.IsAvailable, req.ItemType, req.MenuTypes,
		req.DietaryTags, req.Allergens, req.IsAlcoholic,
	)

	return h.scanMenuItemRowWithoutCategory(row)
}

// Update updates an existing menu item
func (h *DBHandler) Update(id string, req *models.MenuItemUpdateRequest) (*models.MenuItem, error) {
	query, err := h.queries.Get(menuItemSQL.UpdateMenuItemQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id,
		req.Name, req.Description, req.CategoryID, req.Price, req.HappyHourPrice,
		req.ImageURL, req.IsAvailable, req.ItemType, req.MenuTypes,
		req.DietaryTags, req.Allergens, req.IsAlcoholic,
	)

	return h.scanMenuItemRowWithoutCategory(row)
}

// Delete deletes a menu item
func (h *DBHandler) Delete(id string) error {
	deleteQuery, err := h.queries.Get(menuItemSQL.DeleteMenuItemQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete menu item: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("menu item not found")
	}

	h.logger.WithField("id", id).Info("Menu item deleted")
	return nil
}

// UpdateAvailability updates the availability of a menu item
func (h *DBHandler) UpdateAvailability(id string, isAvailable bool) (*models.MenuItem, error) {
	query, err := h.queries.Get(menuItemSQL.UpdateMenuItemAvailabilityQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id, isAvailable)
	return h.scanMenuItemRowWithoutCategory(row)
}

// UpdateImage updates the image URL of a menu item
func (h *DBHandler) UpdateImage(id string, imageURL string) (*models.MenuItem, error) {
	query, err := h.queries.Get(menuItemSQL.UpdateMenuItemImageQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id, imageURL)
	return h.scanMenuItemRowWithoutCategory(row)
}

// UpdateCost updates the item cost
func (h *DBHandler) UpdateCost(id string, cost float64) (*models.MenuItem, error) {
	query, err := h.queries.Get(menuItemSQL.UpdateMenuItemCostQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id, cost)
	return h.scanMenuItemRowWithoutCategory(row)
}

// Helper functions for scanning
func (h *DBHandler) scanMenuItem(rows *sql.Rows) (*models.MenuItem, error) {
	var item models.MenuItem
	var description, categoryName, imageURL sql.NullString
	var itemCost, happyHourPrice sql.NullFloat64
	var dietaryTags, allergens []byte

	err := rows.Scan(
		&item.ID, &item.Name, &description, &item.CategoryID, &categoryName,
		&item.Price, &itemCost, &happyHourPrice, &imageURL, &item.IsAvailable,
		&item.ItemType, &item.MenuTypes, &dietaryTags, &allergens, &item.IsAlcoholic,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan menu item: %w", err)
	}

	if description.Valid {
		item.Description = &description.String
	}
	if categoryName.Valid {
		item.CategoryName = categoryName.String
	}
	if itemCost.Valid {
		item.ItemCost = &itemCost.Float64
	}
	if happyHourPrice.Valid {
		item.HappyHourPrice = &happyHourPrice.Float64
	}
	if imageURL.Valid {
		item.ImageURL = &imageURL.String
	}
	if dietaryTags != nil {
		item.DietaryTags = json.RawMessage(dietaryTags)
	}
	if allergens != nil {
		item.Allergens = json.RawMessage(allergens)
	}

	return &item, nil
}

func (h *DBHandler) scanMenuItemRow(row *sql.Row) (*models.MenuItem, error) {
	var item models.MenuItem
	var description, categoryName, imageURL sql.NullString
	var itemCost, happyHourPrice sql.NullFloat64
	var dietaryTags, allergens []byte

	err := row.Scan(
		&item.ID, &item.Name, &description, &item.CategoryID, &categoryName,
		&item.Price, &itemCost, &happyHourPrice, &imageURL, &item.IsAvailable,
		&item.ItemType, &item.MenuTypes, &dietaryTags, &allergens, &item.IsAlcoholic,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan menu item: %w", err)
	}

	if description.Valid {
		item.Description = &description.String
	}
	if categoryName.Valid {
		item.CategoryName = categoryName.String
	}
	if itemCost.Valid {
		item.ItemCost = &itemCost.Float64
	}
	if happyHourPrice.Valid {
		item.HappyHourPrice = &happyHourPrice.Float64
	}
	if imageURL.Valid {
		item.ImageURL = &imageURL.String
	}
	if dietaryTags != nil {
		item.DietaryTags = json.RawMessage(dietaryTags)
	}
	if allergens != nil {
		item.Allergens = json.RawMessage(allergens)
	}

	return &item, nil
}

func (h *DBHandler) scanMenuItemRowWithoutCategory(row *sql.Row) (*models.MenuItem, error) {
	var item models.MenuItem
	var description, imageURL sql.NullString
	var itemCost, happyHourPrice sql.NullFloat64
	var dietaryTags, allergens []byte

	err := row.Scan(
		&item.ID, &item.Name, &description, &item.CategoryID,
		&item.Price, &itemCost, &happyHourPrice, &imageURL, &item.IsAvailable,
		&item.ItemType, &item.MenuTypes, &dietaryTags, &allergens, &item.IsAlcoholic,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan menu item: %w", err)
	}

	if description.Valid {
		item.Description = &description.String
	}
	if itemCost.Valid {
		item.ItemCost = &itemCost.Float64
	}
	if happyHourPrice.Valid {
		item.HappyHourPrice = &happyHourPrice.Float64
	}
	if imageURL.Valid {
		item.ImageURL = &imageURL.String
	}
	if dietaryTags != nil {
		item.DietaryTags = json.RawMessage(dietaryTags)
	}
	if allergens != nil {
		item.Allergens = json.RawMessage(allergens)
	}

	return &item, nil
}
