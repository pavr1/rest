package handlers

import (
	"context"
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
	queries, err := menuIngredientSQL.LoadQueries(logger)
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
		var notes, stockVariantID, stockVariantName, menuSubCategoryID, menuSubCategoryName sql.NullString

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.MenuVariantID,
			&stockVariantID,
			&stockVariantName,
			&menuSubCategoryID,
			&menuSubCategoryName,
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
		if stockVariantID.Valid {
			ingredient.StockVariantID = &stockVariantID.String
			ingredient.StockVariantName = stockVariantName.String
		}
		if menuSubCategoryID.Valid {
			ingredient.MenuSubCategoryID = &menuSubCategoryID.String
			ingredient.MenuSubCategoryName = menuSubCategoryName.String
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
	var notes, stockVariantID, stockVariantName, menuSubCategoryID, menuSubCategoryName sql.NullString

	err = h.db.QueryRow(query, id).Scan(
		&ingredient.ID,
		&ingredient.MenuVariantID,
		&stockVariantID,
		&stockVariantName,
		&menuSubCategoryID,
		&menuSubCategoryName,
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
	if stockVariantID.Valid {
		ingredient.StockVariantID = &stockVariantID.String
		ingredient.StockVariantName = stockVariantName.String
	}
	if menuSubCategoryID.Valid {
		ingredient.MenuSubCategoryID = &menuSubCategoryID.String
		ingredient.MenuSubCategoryName = menuSubCategoryName.String
	}

	return &ingredient, nil
}

// Create creates a new menu ingredient and updates menu_variant.item_cost in the same transaction
func (h *DBHandler) Create(req models.MenuIngredientCreateRequest, menuVariantID string) (*models.MenuIngredient, error) {
	tx, err := h.db.BeginTx(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query, err := h.queries.Get("create_menu_ingredient")
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ingredient models.MenuIngredient
	var notes, stockVariantID, menuSubCategoryID sql.NullString

	err = tx.QueryRow(query, menuVariantID, req.StockVariantID, req.MenuSubCategoryID, req.Quantity, req.IsOptional, req.Notes).Scan(
		&ingredient.ID,
		&ingredient.MenuVariantID,
		&stockVariantID,
		&menuSubCategoryID,
		&ingredient.Quantity,
		&ingredient.IsOptional,
		&notes,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("menu ingredient already exists for this menu variant")
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23514" {
			return nil, fmt.Errorf("must specify either stock_variant_id or menu_sub_category_id, but not both")
		}
		return nil, fmt.Errorf("failed to create menu ingredient: %w", err)
	}

	if notes.Valid {
		ingredient.Notes = &notes.String
	}
	if stockVariantID.Valid {
		ingredient.StockVariantID = &stockVariantID.String
	}
	if menuSubCategoryID.Valid {
		ingredient.MenuSubCategoryID = &menuSubCategoryID.String
	}

	if err := h.updateMenuVariantItemCostTx(tx, menuVariantID); err != nil {
		return nil, fmt.Errorf("failed to update menu variant item cost: %w", err)
	}
	if err := h.updateMenuSubCategoryItemCostByVariantTx(tx, menuVariantID); err != nil {
		return nil, fmt.Errorf("failed to update menu sub-category item cost: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &ingredient, nil
}

// Update updates a menu ingredient and updates menu_variant.item_cost in the same transaction
func (h *DBHandler) Update(id string, req models.MenuIngredientUpdateRequest) (*models.MenuIngredient, error) {
	tx, err := h.db.BeginTx(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query, err := h.queries.Get("update_menu_ingredient")
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ingredient models.MenuIngredient
	var notes, stockVariantID, menuSubCategoryID sql.NullString

	err = tx.QueryRow(query, id, req.Quantity, req.IsOptional, req.Notes).Scan(
		&ingredient.ID,
		&ingredient.MenuVariantID,
		&stockVariantID,
		&menuSubCategoryID,
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
	if stockVariantID.Valid {
		ingredient.StockVariantID = &stockVariantID.String
	}
	if menuSubCategoryID.Valid {
		ingredient.MenuSubCategoryID = &menuSubCategoryID.String
	}

	if err := h.updateMenuVariantItemCostTx(tx, ingredient.MenuVariantID); err != nil {
		return nil, fmt.Errorf("failed to update menu variant item cost: %w", err)
	}
	if err := h.updateMenuSubCategoryItemCostByVariantTx(tx, ingredient.MenuVariantID); err != nil {
		return nil, fmt.Errorf("failed to update menu sub-category item cost: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &ingredient, nil
}

// Delete deletes a menu ingredient and updates menu_variant.item_cost in the same transaction
func (h *DBHandler) Delete(id string) error {
	tx, err := h.db.BeginTx(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query, err := h.queries.Get("delete_menu_ingredient")
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	var menuVariantID string
	err = tx.QueryRow(query, id).Scan(&menuVariantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("menu ingredient not found")
		}
		return fmt.Errorf("failed to delete menu ingredient: %w", err)
	}

	if err := h.updateMenuVariantItemCostTx(tx, menuVariantID); err != nil {
		return fmt.Errorf("failed to update menu variant item cost: %w", err)
	}
	if err := h.updateMenuSubCategoryItemCostByVariantTx(tx, menuVariantID); err != nil {
		return fmt.Errorf("failed to update menu sub-category item cost: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// updateMenuVariantItemCostTx recalculates menu_variant.item_cost from ingredients (quantity * stock_variants.avg_cost) within tx
func (h *DBHandler) updateMenuVariantItemCostTx(tx *sql.Tx, menuVariantID string) error {
	query, err := h.queries.Get(menuIngredientSQL.RecalculateMenuVariantItemCostQuery)
	if err != nil {
		return fmt.Errorf("failed to get recalculate query: %w", err)
	}
	_, err = tx.Exec(query, menuVariantID)
	return err
}

// updateMenuSubCategoryItemCostByVariantTx recalculates menu_sub_categories.item_cost (AVG of variants in that sub_category) within tx
func (h *DBHandler) updateMenuSubCategoryItemCostByVariantTx(tx *sql.Tx, menuVariantID string) error {
	query, err := h.queries.Get(menuIngredientSQL.RecalculateMenuSubCategoryItemCostByVariantQuery)
	if err != nil {
		return fmt.Errorf("failed to get recalculate sub-category cost query: %w", err)
	}
	_, err = tx.Exec(query, menuVariantID)
	return err
}

// RecalculateCostsByStockVariantID recalculates menu_variant.item_cost and menu_sub_category.item_cost for all menu variants that use the given stock variant (e.g. after outcome invoice updates avg_cost).
func (h *DBHandler) RecalculateCostsByStockVariantID(stockVariantID string) error {
	query, err := h.queries.Get(menuIngredientSQL.GetMenuVariantIDsByStockVariantQuery)
	if err != nil {
		return fmt.Errorf("failed to get menu variant ids query: %w", err)
	}
	rows, err := h.db.Query(query, stockVariantID)
	if err != nil {
		h.logger.WithError(err).Error("failed to list menu variants by stock variant")
		return fmt.Errorf("failed to list menu variants by stock variant: %w", err)
	}
	defer rows.Close()

	var variantIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			h.logger.WithError(err).Error("failed to scan menu variant id")
			return fmt.Errorf("failed to scan menu variant id: %w", err)
		}
		variantIDs = append(variantIDs, id)
	}
	if err := rows.Err(); err != nil {
		h.logger.WithError(err).Error("failed to iterate menu variant ids")
		return fmt.Errorf("error iterating menu variant ids: %w", err)
	}

	if len(variantIDs) == 0 {
		h.logger.WithField("stock_variant_id", stockVariantID).Info("no menu variants found for stock variant")
		return nil
	}

	tx, err := h.db.BeginTx(context.Background())
	if err != nil {
		h.logger.WithError(err).Error("failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, menuVariantID := range variantIDs {
		if err := h.updateMenuVariantItemCostTx(tx, menuVariantID); err != nil {
			return fmt.Errorf("failed to update menu variant item cost for %s: %w", menuVariantID, err)
		}
		if err := h.updateMenuSubCategoryItemCostByVariantTx(tx, menuVariantID); err != nil {
			return fmt.Errorf("failed to update menu sub-category item cost for variant %s: %w", menuVariantID, err)
		}
	}

	return tx.Commit()
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
		var notes, stockVariantID, stockVariantName, menuSubCategoryID, menuSubCategoryName sql.NullString
		var stockVariantAvgCost sql.NullFloat64

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.MenuVariantID,
			&stockVariantID,
			&stockVariantName,
			&stockVariantAvgCost,
			&menuSubCategoryID,
			&menuSubCategoryName,
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
		if stockVariantID.Valid {
			ingredient.StockVariantID = &stockVariantID.String
			ingredient.StockVariantName = stockVariantName.String
			if stockVariantAvgCost.Valid {
				ingredient.StockVariantAvgCost = &stockVariantAvgCost.Float64
			}
		}
		if menuSubCategoryID.Valid {
			ingredient.MenuSubCategoryID = &menuSubCategoryID.String
			ingredient.MenuSubCategoryName = menuSubCategoryName.String
		}

		ingredients = append(ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating menu ingredients: %w", err)
	}

	return ingredients, nil
}
