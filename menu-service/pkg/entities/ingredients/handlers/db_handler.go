package handlers

import (
	"database/sql"
	"fmt"
	"menu-service/pkg/entities/ingredients/models"
	ingredientSQL "menu-service/pkg/entities/ingredients/sql"
	sharedDb "shared/db"
	"time"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for ingredients
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *ingredientSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := ingredientSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns all ingredients for a menu item
func (h *DBHandler) List(menuItemID string) (*models.IngredientListResponse, error) {
	query, err := h.queries.Get(ingredientSQL.ListIngredientsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	rows, err := h.db.Query(query, menuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to list ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []models.Ingredient
	for rows.Next() {
		var ing models.Ingredient
		if err := rows.Scan(&ing.ID, &ing.MenuItemID, &ing.StockItemID, &ing.StockItemName,
			&ing.StockItemUnit, &ing.Quantity, &ing.CreatedAt, &ing.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}
		ingredients = append(ingredients, ing)
	}

	return &models.IngredientListResponse{
		MenuItemID:  menuItemID,
		Ingredients: ingredients,
	}, nil
}

// Get returns a specific ingredient
func (h *DBHandler) Get(menuItemID, stockItemID string) (*models.Ingredient, error) {
	query, err := h.queries.Get(ingredientSQL.GetIngredientQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ing models.Ingredient
	err = h.db.QueryRow(query, menuItemID, stockItemID).Scan(
		&ing.ID, &ing.MenuItemID, &ing.StockItemID, &ing.StockItemName,
		&ing.StockItemUnit, &ing.Quantity, &ing.CreatedAt, &ing.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ingredient: %w", err)
	}

	return &ing, nil
}

// Create adds a new ingredient to a menu item
func (h *DBHandler) Create(menuItemID string, req *models.IngredientCreateRequest) (*models.Ingredient, error) {
	query, err := h.queries.Get(ingredientSQL.CreateIngredientQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ing models.Ingredient
	err = h.db.QueryRow(query, menuItemID, req.StockItemID, req.Quantity).Scan(
		&ing.ID, &ing.MenuItemID, &ing.StockItemID, &ing.Quantity, &ing.CreatedAt, &ing.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ingredient: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"menu_item_id":  menuItemID,
		"stock_item_id": req.StockItemID,
	}).Info("Ingredient added to menu item")

	return &ing, nil
}

// Update updates the quantity of an ingredient
func (h *DBHandler) Update(menuItemID, stockItemID string, req *models.IngredientUpdateRequest) (*models.Ingredient, error) {
	query, err := h.queries.Get(ingredientSQL.UpdateIngredientQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var ing models.Ingredient
	err = h.db.QueryRow(query, menuItemID, stockItemID, req.Quantity).Scan(
		&ing.ID, &ing.MenuItemID, &ing.StockItemID, &ing.Quantity, &ing.CreatedAt, &ing.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update ingredient: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"menu_item_id":  menuItemID,
		"stock_item_id": stockItemID,
	}).Info("Ingredient quantity updated")

	return &ing, nil
}

// Delete removes an ingredient from a menu item
func (h *DBHandler) Delete(menuItemID, stockItemID string) error {
	query, err := h.queries.Get(ingredientSQL.DeleteIngredientQuery)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	result, err := h.db.Exec(query, menuItemID, stockItemID)
	if err != nil {
		return fmt.Errorf("failed to delete ingredient: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("ingredient not found")
	}

	h.logger.WithFields(logrus.Fields{
		"menu_item_id":  menuItemID,
		"stock_item_id": stockItemID,
	}).Info("Ingredient removed from menu item")

	return nil
}

// CalculateCost calculates the cost of a menu item based on its ingredients
func (h *DBHandler) CalculateCost(menuItemID, menuItemName string) (*models.MenuItemCost, error) {
	query, err := h.queries.Get(ingredientSQL.CalculateMenuItemCostQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	rows, err := h.db.Query(query, menuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate cost: %w", err)
	}
	defer rows.Close()

	var ingredientCosts []models.IngredientCost
	var totalCost float64

	for rows.Next() {
		var ic models.IngredientCost
		if err := rows.Scan(&ic.StockItemID, &ic.StockItemName, &ic.Quantity, &ic.UnitCost); err != nil {
			return nil, fmt.Errorf("failed to scan cost row: %w", err)
		}
		ic.TotalCost = ic.Quantity * ic.UnitCost
		totalCost += ic.TotalCost
		ingredientCosts = append(ingredientCosts, ic)
	}

	return &models.MenuItemCost{
		MenuItemID:     menuItemID,
		MenuItemName:   menuItemName,
		TotalCost:      totalCost,
		IngredientCost: ingredientCosts,
		CalculatedAt:   time.Now(),
	}, nil
}
