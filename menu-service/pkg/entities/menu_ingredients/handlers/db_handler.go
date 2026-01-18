package handlers

import (
	"database/sql"
	"fmt"
	"menu-service/pkg/entities/menu_ingredients/models"
	menuIngredientSQL "menu-service/pkg/entities/menu_ingredients/sql"
	sharedDb "shared/db"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for menu ingredients
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *menuIngredientSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new menu ingredient database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := menuIngredientSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List retrieves menu ingredients with pagination
func (h *DBHandler) List(page, limit int) ([]models.MenuIngredient, error) {
	offset := (page - 1) * limit

	query, err := h.queries.Get("list_menu_ingredients")
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	rows, err := h.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []models.MenuIngredient
	for rows.Next() {
		var ingredient models.MenuIngredient
		var notes sql.NullString

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.MenuVariantID,
			&ingredient.StockSubCategoryID,
			&ingredient.StockSubCategoryName,
			&ingredient.Quantity,
			&ingredient.IsOptional,
			&notes,
			&ingredient.CreatedAt,
			&ingredient.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu ingredient: %w", err)
		}

		if notes.Valid {
			ingredient.Notes = &notes.String
		}

		ingredients = append(ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating menu ingredients: %w", err)
	}

	return ingredients, nil
}

// GetByID retrieves a menu ingredient by ID
func (h *DBHandler) GetByID(id string) (*models.MenuIngredient, error) {
	query, err := h.queries.Get("get_menu_ingredient_by_id")
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ingredient models.MenuIngredient
	var notes sql.NullString

	err = h.db.QueryRow(query, id).Scan(
		&ingredient.ID,
		&ingredient.MenuVariantID,
		&ingredient.StockSubCategoryID,
		&ingredient.StockSubCategoryName,
		&ingredient.Quantity,
		&ingredient.IsOptional,
		&notes,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get menu ingredient by ID: %w", err)
	}

	if notes.Valid {
		ingredient.Notes = &notes.String
	}

	return &ingredient, nil
}

// Create creates a new menu ingredient
func (h *DBHandler) Create(req models.MenuIngredientCreateRequest, menuVariantID string) (*models.MenuIngredient, error) {
	query, err := h.queries.Get("create_menu_ingredient")
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ingredient models.MenuIngredient
	var notes sql.NullString

	err = h.db.QueryRow(query, menuVariantID, req.StockSubCategoryID, req.Quantity, req.IsOptional, req.Notes).Scan(
		&ingredient.ID,
		&ingredient.MenuVariantID,
		&ingredient.StockSubCategoryID,
		&ingredient.Quantity,
		&ingredient.IsOptional,
		&notes,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("menu ingredient already exists for this menu variant and stock sub-category")
		}
		return nil, fmt.Errorf("failed to create menu ingredient: %w", err)
	}

	if notes.Valid {
		ingredient.Notes = &notes.String
	}

	return &ingredient, nil
}

// Update updates a menu ingredient
func (h *DBHandler) Update(id string, req models.MenuIngredientUpdateRequest) (*models.MenuIngredient, error) {
	query, err := h.queries.Get("update_menu_ingredient")
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ingredient models.MenuIngredient
	var notes sql.NullString

	err = h.db.QueryRow(query, id, req.Quantity, req.IsOptional, req.Notes).Scan(
		&ingredient.ID,
		&ingredient.MenuVariantID,
		&ingredient.StockSubCategoryID,
		&ingredient.Quantity,
		&ingredient.IsOptional,
		&notes,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("menu ingredient not found")
		}
		return nil, fmt.Errorf("failed to update menu ingredient: %w", err)
	}

	if notes.Valid {
		ingredient.Notes = &notes.String
	}

	return &ingredient, nil
}

// Delete deletes a menu ingredient
func (h *DBHandler) Delete(id string) error {
	query, err := h.queries.Get("delete_menu_ingredient")
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	result, err := h.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete menu ingredient: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("menu ingredient not found")
	}

	return nil
}

// GetByMenuVariant retrieves all ingredients for a specific menu variant
func (h *DBHandler) GetByMenuVariant(menuVariantID string) ([]models.MenuIngredient, error) {
	query, err := h.queries.Get("get_ingredients_by_menu_variant")
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	rows, err := h.db.Query(query, menuVariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients by menu variant: %w", err)
	}
	defer rows.Close()

	var ingredients []models.MenuIngredient
	for rows.Next() {
		var ingredient models.MenuIngredient
		var notes sql.NullString

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.MenuVariantID,
			&ingredient.StockSubCategoryID,
			&ingredient.StockSubCategoryName,
			&ingredient.Quantity,
			&ingredient.IsOptional,
			&notes,
			&ingredient.CreatedAt,
			&ingredient.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu ingredient: %w", err)
		}

		if notes.Valid {
			ingredient.Notes = &notes.String
		}

		ingredients = append(ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating menu ingredients: %w", err)
	}

	return ingredients, nil
}
