package handlers

import (
	"database/sql"
	"fmt"
	"menu-service/pkg/entities/stock_item_categories/models"
	stockCategorySQL "menu-service/pkg/entities/stock_item_categories/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for stock item categories
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *stockCategorySQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := stockCategorySQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of stock item categories
func (h *DBHandler) List(page, limit int) (*models.StockItemCategoryListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockCategorySQL.CountStockItemCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock item categories: %w", err)
	}

	listQuery, err := h.queries.Get(stockCategorySQL.ListStockItemCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock item categories: %w", err)
	}
	defer rows.Close()

	var categories []models.StockItemCategory
	for rows.Next() {
		var cat models.StockItemCategory
		var description sql.NullString

		if err := rows.Scan(&cat.ID, &cat.Name, &description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock item category: %w", err)
		}

		if description.Valid {
			cat.Description = &description.String
		}
		categories = append(categories, cat)
	}

	return &models.StockItemCategoryListResponse{
		Categories: categories,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetByID returns a stock item category by ID
func (h *DBHandler) GetByID(id string) (*models.StockItemCategory, error) {
	query, err := h.queries.Get(stockCategorySQL.GetStockItemCategoryByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.StockItemCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &description, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock item category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	return &cat, nil
}

// Create creates a new stock item category
func (h *DBHandler) Create(req *models.StockItemCategoryCreateRequest) (*models.StockItemCategory, error) {
	query, err := h.queries.Get(stockCategorySQL.CreateStockItemCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.StockItemCategory
	var description sql.NullString

	err = h.db.QueryRow(query, req.Name, req.Description).Scan(
		&cat.ID, &cat.Name, &description, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock item category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	h.logger.WithField("id", cat.ID).Info("Stock item category created")
	return &cat, nil
}

// Update updates an existing stock item category
func (h *DBHandler) Update(id string, req *models.StockItemCategoryUpdateRequest) (*models.StockItemCategory, error) {
	query, err := h.queries.Get(stockCategorySQL.UpdateStockItemCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.StockItemCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id, req.Name, req.Description).Scan(
		&cat.ID, &cat.Name, &description, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update stock item category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	h.logger.WithField("id", cat.ID).Info("Stock item category updated")
	return &cat, nil
}

// Delete deletes a stock item category
func (h *DBHandler) Delete(id string) error {
	checkQuery, err := h.queries.Get(stockCategorySQL.CheckStockItemCategoryDependenciesQuery)
	if err != nil {
		return fmt.Errorf("failed to get check query: %w", err)
	}

	var count int
	if err := h.db.QueryRow(checkQuery, id).Scan(&count); err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete category: %d stock items depend on it", count)
	}

	deleteQuery, err := h.queries.Get(stockCategorySQL.DeleteStockItemCategoryQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock item category: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock item category not found")
	}

	h.logger.WithField("id", id).Info("Stock item category deleted")
	return nil
}
