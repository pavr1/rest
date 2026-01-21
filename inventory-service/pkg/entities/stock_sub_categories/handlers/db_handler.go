package handlers

import (
	"database/sql"
	"fmt"
	"inventory-service/pkg/entities/stock_sub_categories/models"
	stockSubCategorySQL "inventory-service/pkg/entities/stock_sub_categories/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for stock sub-categories
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *stockSubCategorySQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := stockSubCategorySQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of stock sub-categories
func (h *DBHandler) List(page, limit int) (*models.StockSubCategoryListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockSubCategorySQL.CountStockSubCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock sub-categories: %w", err)
	}

	listQuery, err := h.queries.Get(stockSubCategorySQL.ListStockSubCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock sub-categories: %w", err)
	}
	defer rows.Close()

	var subCategories []models.StockSubCategory
	for rows.Next() {
		var subCat models.StockSubCategory
		var description sql.NullString

		if err := rows.Scan(&subCat.ID, &subCat.Name, &description, &subCat.StockCategoryID, &subCat.DisplayOrder, &subCat.IsActive, &subCat.CreatedAt, &subCat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock sub-category: %w", err)
		}

		if description.Valid {
			subCat.Description = &description.String
		}
		subCategories = append(subCategories, subCat)
	}

	return &models.StockSubCategoryListResponse{
		SubCategories: subCategories,
		Total:         total,
		Page:          page,
		Limit:         limit,
	}, nil
}

// ListByCategory returns a paginated list of stock sub-categories filtered by category
func (h *DBHandler) ListByCategory(categoryID string, page, limit int) (*models.StockSubCategoryListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockSubCategorySQL.CountStockSubCategoriesByCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, categoryID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock sub-categories: %w", err)
	}

	listQuery, err := h.queries.Get(stockSubCategorySQL.ListStockSubCategoriesByCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock sub-categories: %w", err)
	}
	defer rows.Close()

	var subCategories []models.StockSubCategory
	for rows.Next() {
		var subCat models.StockSubCategory
		var description sql.NullString

		if err := rows.Scan(&subCat.ID, &subCat.Name, &description, &subCat.StockCategoryID, &subCat.DisplayOrder, &subCat.IsActive, &subCat.CreatedAt, &subCat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock sub-category: %w", err)
		}

		if description.Valid {
			subCat.Description = &description.String
		}
		subCategories = append(subCategories, subCat)
	}

	return &models.StockSubCategoryListResponse{
		SubCategories: subCategories,
		Total:         total,
		Page:          page,
		Limit:         limit,
	}, nil
}

// GetByID returns a stock sub-category by ID
func (h *DBHandler) GetByID(id string) (*models.StockSubCategory, error) {
	query, err := h.queries.Get(stockSubCategorySQL.GetStockSubCategoryByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var subCat models.StockSubCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id).Scan(&subCat.ID, &subCat.Name, &description, &subCat.StockCategoryID, &subCat.DisplayOrder, &subCat.IsActive, &subCat.CreatedAt, &subCat.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock sub-category: %w", err)
	}

	if description.Valid {
		subCat.Description = &description.String
	}

	return &subCat, nil
}

// Create creates a new stock sub-category
func (h *DBHandler) Create(req *models.StockSubCategoryCreateRequest) (*models.StockSubCategory, error) {
	query, err := h.queries.Get(stockSubCategorySQL.CreateStockSubCategoryQuery)
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

	var subCat models.StockSubCategory
	var description sql.NullString

	err = h.db.QueryRow(query, req.Name, req.Description, req.StockCategoryID, displayOrder, isActive).Scan(
		&subCat.ID, &subCat.Name, &description, &subCat.StockCategoryID, &subCat.DisplayOrder, &subCat.IsActive, &subCat.CreatedAt, &subCat.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock sub-category: %w", err)
	}

	if description.Valid {
		subCat.Description = &description.String
	}

	h.logger.WithField("id", subCat.ID).Info("Stock sub-category created")
	return &subCat, nil
}

// Update updates an existing stock sub-category
func (h *DBHandler) Update(id string, req *models.StockSubCategoryUpdateRequest) (*models.StockSubCategory, error) {
	query, err := h.queries.Get(stockSubCategorySQL.UpdateStockSubCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var subCat models.StockSubCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id, req.Name, req.Description, req.DisplayOrder, req.IsActive).Scan(
		&subCat.ID, &subCat.Name, &description, &subCat.StockCategoryID, &subCat.DisplayOrder, &subCat.IsActive, &subCat.CreatedAt, &subCat.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update stock sub-category: %w", err)
	}

	if description.Valid {
		subCat.Description = &description.String
	}

	h.logger.WithField("id", subCat.ID).Info("Stock sub-category updated")
	return &subCat, nil
}

// Delete deletes a stock sub-category
func (h *DBHandler) Delete(id string) error {
	checkQuery, err := h.queries.Get(stockSubCategorySQL.CheckStockSubCategoryDependenciesQuery)
	if err != nil {
		return fmt.Errorf("failed to get check query: %w", err)
	}

	var count int
	if err := h.db.QueryRow(checkQuery, id).Scan(&count); err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete sub-category: %d stock variants depend on it", count)
	}

	deleteQuery, err := h.queries.Get(stockSubCategorySQL.DeleteStockSubCategoryQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock sub-category: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock sub-category not found")
	}

	h.logger.WithField("id", id).Info("Stock sub-category deleted")
	return nil
}
