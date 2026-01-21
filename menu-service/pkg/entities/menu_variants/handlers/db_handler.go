package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"menu-service/pkg/entities/menu_variants/models"
	menuVariantSQL "menu-service/pkg/entities/menu_variants/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for menu variants
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *menuVariantSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := menuVariantSQL.LoadQueries()
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
func (h *DBHandler) List(req *models.MenuVariantListRequest) (*models.MenuVariantListResponse, error) {
	offset := (req.Page - 1) * req.Limit

	// Prepare menu_types filter as JSONB
	var menuTypesJSON interface{}
	if req.MenuType != nil {
		menuTypesJSON = fmt.Sprintf(`["%s"]`, *req.MenuType)
	}

	countQuery, err := h.queries.Get(menuVariantSQL.CountMenuVariantsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, req.CategoryID, req.SubCategoryID, req.IsAvailable, menuTypesJSON).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count menu items: %w", err)
	}

	listQuery, err := h.queries.Get(menuVariantSQL.ListMenuVariantsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, req.CategoryID, req.SubCategoryID, req.IsAvailable, menuTypesJSON, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items: %w", err)
	}
	defer rows.Close()

	var items []models.MenuVariant
	for rows.Next() {
		item, err := h.scanMenuVariant(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return &models.MenuVariantListResponse{
		Items: items,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// GetByID returns a menu item by ID
func (h *DBHandler) GetByID(id string) (*models.MenuVariant, error) {
	query, err := h.queries.Get(menuVariantSQL.GetMenuVariantByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id)
	return h.scanMenuVariantRow(row)
}

// Create creates a new menu item
func (h *DBHandler) Create(req *models.MenuVariantCreateRequest) (*models.MenuVariant, error) {
	query, err := h.queries.Get(menuVariantSQL.CreateMenuVariantQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	// Default empty JSON arrays for JSONB columns if nil
	menuTypes := req.MenuTypes
	if menuTypes == nil {
		menuTypes = json.RawMessage(`[]`)
	}
	dietaryTags := req.DietaryTags
	if dietaryTags == nil {
		dietaryTags = json.RawMessage(`[]`)
	}
	allergens := req.Allergens
	if allergens == nil {
		allergens = json.RawMessage(`[]`)
	}

	row := h.db.QueryRow(query,
		req.Name, req.Description, req.SubCategoryID, req.Price, req.HappyHourPrice,
		req.ImageURL, req.IsAvailable, req.PreparationTime, menuTypes,
		dietaryTags, allergens, req.IsAlcoholic, req.DisplayOrder,
	)

	return h.scanMenuVariantRowWithoutSubCategory(row)
}

// Update updates an existing menu item
func (h *DBHandler) Update(id string, req *models.MenuVariantUpdateRequest) (*models.MenuVariant, error) {
	query, err := h.queries.Get(menuVariantSQL.UpdateMenuVariantQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	// For update, we pass nil to keep existing values, or the new value
	// The SQL uses COALESCE to handle this
	var menuTypes, dietaryTags, allergens interface{}
	if req.MenuTypes != nil {
		menuTypes = *req.MenuTypes
	}
	if req.DietaryTags != nil {
		dietaryTags = *req.DietaryTags
	}
	if req.Allergens != nil {
		allergens = *req.Allergens
	}

	row := h.db.QueryRow(query, id,
		req.Name, req.Description, req.SubCategoryID, req.Price, req.HappyHourPrice,
		req.ImageURL, req.IsAvailable, req.PreparationTime, menuTypes,
		dietaryTags, allergens, req.IsAlcoholic, req.DisplayOrder,
	)

	return h.scanMenuVariantRowWithoutSubCategory(row)
}

// Delete deletes a menu item
func (h *DBHandler) Delete(id string) error {
	deleteQuery, err := h.queries.Get(menuVariantSQL.DeleteMenuVariantQuery)
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
func (h *DBHandler) UpdateAvailability(id string, isAvailable bool) (*models.MenuVariant, error) {
	query, err := h.queries.Get(menuVariantSQL.UpdateMenuVariantAvailabilityQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id, isAvailable)
	return h.scanMenuVariantRowWithoutSubCategory(row)
}

// UpdateImage updates the image URL of a menu item
func (h *DBHandler) UpdateImage(id string, imageURL string) (*models.MenuVariant, error) {
	query, err := h.queries.Get(menuVariantSQL.UpdateMenuVariantImageQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id, imageURL)
	return h.scanMenuVariantRowWithoutSubCategory(row)
}

// UpdateCost updates the item cost
func (h *DBHandler) UpdateCost(id string, cost float64) (*models.MenuVariant, error) {
	query, err := h.queries.Get(menuVariantSQL.UpdateMenuVariantCostQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id, cost)
	return h.scanMenuVariantRowWithoutSubCategory(row)
}

// Helper functions for scanning
func (h *DBHandler) scanMenuVariant(rows *sql.Rows) (*models.MenuVariant, error) {
	var item models.MenuVariant
	var description, subMenuName, itemType, imageURL sql.NullString
	var itemCost, happyHourPrice sql.NullFloat64
	var preparationTime sql.NullInt32
	var dietaryTags, allergens []byte

	err := rows.Scan(
		&item.ID, &item.Name, &description, &item.SubCategoryID, &subMenuName,
		&itemType, &item.Price, &itemCost, &happyHourPrice, &imageURL, &item.IsAvailable,
		&preparationTime, &item.MenuTypes, &dietaryTags, &allergens, &item.IsAlcoholic,
		&item.DisplayOrder, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan menu item: %w", err)
	}

	if description.Valid {
		item.Description = &description.String
	}
	if subMenuName.Valid {
		item.SubCategoryName = subMenuName.String
	}
	if itemType.Valid {
		item.ItemType = itemType.String
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
	if preparationTime.Valid {
		prepTime := int(preparationTime.Int32)
		item.PreparationTime = &prepTime
	}
	if dietaryTags != nil {
		item.DietaryTags = json.RawMessage(dietaryTags)
	}
	if allergens != nil {
		item.Allergens = json.RawMessage(allergens)
	}

	return &item, nil
}

func (h *DBHandler) scanMenuVariantRow(row *sql.Row) (*models.MenuVariant, error) {
	var item models.MenuVariant
	var description, subMenuName, itemType, imageURL sql.NullString
	var itemCost, happyHourPrice sql.NullFloat64
	var preparationTime sql.NullInt32
	var dietaryTags, allergens []byte

	err := row.Scan(
		&item.ID, &item.Name, &description, &item.SubCategoryID, &subMenuName,
		&itemType, &item.Price, &itemCost, &happyHourPrice, &imageURL, &item.IsAvailable,
		&preparationTime, &item.MenuTypes, &dietaryTags, &allergens, &item.IsAlcoholic,
		&item.DisplayOrder, &item.CreatedAt, &item.UpdatedAt,
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
	if subMenuName.Valid {
		item.SubCategoryName = subMenuName.String
	}
	if itemType.Valid {
		item.ItemType = itemType.String
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
	if preparationTime.Valid {
		prepTime := int(preparationTime.Int32)
		item.PreparationTime = &prepTime
	}
	if dietaryTags != nil {
		item.DietaryTags = json.RawMessage(dietaryTags)
	}
	if allergens != nil {
		item.Allergens = json.RawMessage(allergens)
	}

	return &item, nil
}

func (h *DBHandler) scanMenuVariantRowWithoutSubCategory(row *sql.Row) (*models.MenuVariant, error) {
	var item models.MenuVariant
	var description, imageURL sql.NullString
	var itemCost, happyHourPrice sql.NullFloat64
	var preparationTime sql.NullInt32
	var dietaryTags, allergens []byte

	err := row.Scan(
		&item.ID, &item.Name, &description, &item.SubCategoryID,
		&item.Price, &itemCost, &happyHourPrice, &imageURL, &item.IsAvailable,
		&preparationTime, &item.MenuTypes, &dietaryTags, &allergens, &item.IsAlcoholic,
		&item.DisplayOrder, &item.CreatedAt, &item.UpdatedAt,
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
	if preparationTime.Valid {
		prepTime := int(preparationTime.Int32)
		item.PreparationTime = &prepTime
	}
	if dietaryTags != nil {
		item.DietaryTags = json.RawMessage(dietaryTags)
	}
	if allergens != nil {
		item.Allergens = json.RawMessage(allergens)
	}

	return &item, nil
}
