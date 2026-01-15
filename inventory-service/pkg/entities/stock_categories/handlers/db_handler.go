package handlers

import (
	"database/sql"
	"fmt"
	"inventory-service/pkg/entities/stock_categories/models"
	stockCategorySQL "inventory-service/pkg/entities/stock_categories/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for stock categories
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

// List returns a paginated list of stock categories
func (h *DBHandler) List(page, limit int) (*models.StockCategoryListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockCategorySQL.CountStockCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock categories: %w", err)
	}

	listQuery, err := h.queries.Get(stockCategorySQL.ListStockCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock categories: %w", err)
	}
	defer rows.Close()

	var categories []models.StockCategory
	for rows.Next() {
		var cat models.StockCategory
		var description sql.NullString

		if err := rows.Scan(&cat.ID, &cat.Name, &description, &cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock category: %w", err)
		}

		if description.Valid {
			cat.Description = &description.String
		}
		categories = append(categories, cat)
	}

	return &models.StockCategoryListResponse{
		Categories: categories,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetByID returns a stock category by ID
func (h *DBHandler) GetByID(id string) (*models.StockCategory, error) {
	query, err := h.queries.Get(stockCategorySQL.GetStockCategoryByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.StockCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &description, &cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	return &cat, nil
}

// Create creates a new stock category
func (h *DBHandler) Create(req *models.StockCategoryCreateRequest) (*models.StockCategory, error) {
	query, err := h.queries.Get(stockCategorySQL.CreateStockCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	// Set defaults if not provided
	displayOrder := 0
	if req.DisplayOrder != nil {
		displayOrder = *req.DisplayOrder
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var cat models.StockCategory
	var description sql.NullString

	err = h.db.QueryRow(query, req.Name, req.Description, displayOrder, isActive).Scan(
		&cat.ID, &cat.Name, &description, &cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	h.logger.WithField("id", cat.ID).Info("Stock category created")
	return &cat, nil
}

// Update updates an existing stock category
func (h *DBHandler) Update(id string, req *models.StockCategoryUpdateRequest) (*models.StockCategory, error) {
	query, err := h.queries.Get(stockCategorySQL.UpdateStockCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.StockCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id, req.Name, req.Description, req.DisplayOrder, req.IsActive).Scan(
		&cat.ID, &cat.Name, &description, &cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update stock category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	h.logger.WithField("id", cat.ID).Info("Stock category updated")
	return &cat, nil
}

// Delete deletes a stock category
func (h *DBHandler) Delete(id string) error {
	checkQuery, err := h.queries.Get(stockCategorySQL.CheckStockCategoryDependenciesQuery)
	if err != nil {
		return fmt.Errorf("failed to get check query: %w", err)
	}

	var count int
	if err := h.db.QueryRow(checkQuery, id).Scan(&count); err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete category: %d stock sub-categories depend on it", count)
	}

	deleteQuery, err := h.queries.Get(stockCategorySQL.DeleteStockCategoryQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock category: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock category not found")
	}

	h.logger.WithField("id", id).Info("Stock category deleted")
	return nil
}
