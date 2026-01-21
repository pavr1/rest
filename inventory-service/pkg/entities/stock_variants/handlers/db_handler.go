package handlers

import (
	"database/sql"
	"fmt"
	"inventory-service/pkg/entities/stock_variants/models"
	stockVariantSQL "inventory-service/pkg/entities/stock_variants/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for stock variants
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *stockVariantSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := stockVariantSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of stock variants
func (h *DBHandler) List(page, limit int) (*models.StockVariantListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockVariantSQL.CountStockVariantsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock variants: %w", err)
	}

	listQuery, err := h.queries.Get(stockVariantSQL.ListStockVariantsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock variants: %w", err)
	}
	defer rows.Close()

	var variants []models.StockVariant
	for rows.Next() {
		var variant models.StockVariant

		if err := rows.Scan(&variant.ID, &variant.Name, &variant.Description, &variant.StockSubCategoryID, &variant.AvgCost, &variant.IsActive, &variant.CreatedAt, &variant.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock variant: %w", err)
		}

		variants = append(variants, variant)
	}

	return &models.StockVariantListResponse{
		Variants: variants,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// ListByCategory returns a paginated list of stock variants filtered by category
func (h *DBHandler) ListByCategory(categoryID string, page, limit int) (*models.StockVariantListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockVariantSQL.CountStockVariantsByCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, categoryID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock variants: %w", err)
	}

	listQuery, err := h.queries.Get(stockVariantSQL.ListStockVariantsByCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock variants: %w", err)
	}
	defer rows.Close()

	var variants []models.StockVariant
	for rows.Next() {
		var variant models.StockVariant

		if err := rows.Scan(&variant.ID, &variant.Name, &variant.Description, &variant.StockSubCategoryID, &variant.AvgCost, &variant.IsActive, &variant.CreatedAt, &variant.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock variant: %w", err)
		}

		variants = append(variants, variant)
	}

	return &models.StockVariantListResponse{
		Variants: variants,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// ListBySubCategory returns a paginated list of stock variants filtered by sub-category
func (h *DBHandler) ListBySubCategory(subCategoryID string, page, limit int) (*models.StockVariantListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockVariantSQL.CountStockVariantsBySubCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, subCategoryID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock variants: %w", err)
	}

	listQuery, err := h.queries.Get(stockVariantSQL.ListStockVariantsBySubCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, subCategoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock variants: %w", err)
	}
	defer rows.Close()

	var variants []models.StockVariant
	for rows.Next() {
		var variant models.StockVariant

		if err := rows.Scan(&variant.ID, &variant.Name, &variant.Description, &variant.StockSubCategoryID, &variant.AvgCost, &variant.IsActive, &variant.CreatedAt, &variant.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock variant: %w", err)
		}

		variants = append(variants, variant)
	}

	return &models.StockVariantListResponse{
		Variants: variants,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// GetByID returns a stock variant by ID
func (h *DBHandler) GetByID(id string) (*models.StockVariant, error) {
	query, err := h.queries.Get(stockVariantSQL.GetStockVariantByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var variant models.StockVariant

	err = h.db.QueryRow(query, id).Scan(&variant.ID, &variant.Name, &variant.Description, &variant.StockSubCategoryID, &variant.AvgCost, &variant.IsActive, &variant.CreatedAt, &variant.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock variant: %w", err)
	}

	return &variant, nil
}

// Create creates a new stock variant
func (h *DBHandler) Create(req *models.StockVariantCreateRequest) (*models.StockVariant, error) {
	query, err := h.queries.Get(stockVariantSQL.CreateStockVariantQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	// Set defaults if not provided
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var variant models.StockVariant

	err = h.db.QueryRow(query, req.Name, req.Description, req.StockSubCategoryID, isActive).Scan(
		&variant.ID, &variant.Name, &variant.Description, &variant.StockSubCategoryID, &variant.AvgCost, &variant.IsActive, &variant.CreatedAt, &variant.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock variant: %w", err)
	}

	h.logger.WithField("id", variant.ID).Info("Stock variant created")
	return &variant, nil
}

// Update updates an existing stock variant
func (h *DBHandler) Update(id string, req *models.StockVariantUpdateRequest) (*models.StockVariant, error) {
	query, err := h.queries.Get(stockVariantSQL.UpdateStockVariantQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var variant models.StockVariant

	err = h.db.QueryRow(query, id, req.Name, req.Description, req.IsActive).Scan(
		&variant.ID, &variant.Name, &variant.Description, &variant.StockSubCategoryID, &variant.AvgCost, &variant.IsActive, &variant.CreatedAt, &variant.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update stock variant: %w", err)
	}

	h.logger.WithField("id", variant.ID).Info("Stock variant updated")
	return &variant, nil
}

// Delete deletes a stock variant
func (h *DBHandler) Delete(id string) error {
	checkQuery, err := h.queries.Get(stockVariantSQL.CheckStockVariantDependenciesQuery)
	if err != nil {
		return fmt.Errorf("failed to get check query: %w", err)
	}

	var count int
	if err := h.db.QueryRow(checkQuery, id).Scan(&count); err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete variant: %d dependencies found", count)
	}

	deleteQuery, err := h.queries.Get(stockVariantSQL.DeleteStockVariantQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock variant: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock variant not found")
	}

	h.logger.WithField("id", id).Info("Stock variant deleted")
	return nil
}
