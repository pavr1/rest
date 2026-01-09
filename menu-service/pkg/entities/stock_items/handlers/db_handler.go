package handlers

import (
	"database/sql"
	"fmt"
	"menu-service/pkg/entities/stock_items/models"
	stockItemSQL "menu-service/pkg/entities/stock_items/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for stock items
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *stockItemSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := stockItemSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of stock items
func (h *DBHandler) List(req *models.StockItemListRequest) (*models.StockItemListResponse, error) {
	offset := (req.Page - 1) * req.Limit

	countQuery, err := h.queries.Get(stockItemSQL.CountStockItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, req.CategoryID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock items: %w", err)
	}

	listQuery, err := h.queries.Get(stockItemSQL.ListStockItemsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, req.CategoryID, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock items: %w", err)
	}
	defer rows.Close()

	var items []models.StockItem
	for rows.Next() {
		var item models.StockItem
		var description, categoryID, categoryName sql.NullString

		if err := rows.Scan(&item.ID, &item.Name, &item.Unit, &description, &categoryID, &categoryName, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock item: %w", err)
		}

		if description.Valid {
			item.Description = &description.String
		}
		if categoryID.Valid {
			item.CategoryID = &categoryID.String
		}
		if categoryName.Valid {
			item.CategoryName = &categoryName.String
		}
		items = append(items, item)
	}

	return &models.StockItemListResponse{
		Items: items,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// GetByID returns a stock item by ID
func (h *DBHandler) GetByID(id string) (*models.StockItem, error) {
	query, err := h.queries.Get(stockItemSQL.GetStockItemByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var item models.StockItem
	var description, categoryID, categoryName sql.NullString

	err = h.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Unit, &description, &categoryID, &categoryName, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock item: %w", err)
	}

	if description.Valid {
		item.Description = &description.String
	}
	if categoryID.Valid {
		item.CategoryID = &categoryID.String
	}
	if categoryName.Valid {
		item.CategoryName = &categoryName.String
	}

	return &item, nil
}

// Create creates a new stock item
func (h *DBHandler) Create(req *models.StockItemCreateRequest) (*models.StockItem, error) {
	// Validate unit
	if !models.IsValidUnit(req.Unit) {
		return nil, fmt.Errorf("invalid unit: %s. Valid units are: %v", req.Unit, models.ValidUnits)
	}

	query, err := h.queries.Get(stockItemSQL.CreateStockItemQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var item models.StockItem
	var description, categoryID sql.NullString

	err = h.db.QueryRow(query, req.Name, req.Unit, req.Description, req.CategoryID).Scan(
		&item.ID, &item.Name, &item.Unit, &description, &categoryID, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock item: %w", err)
	}

	if description.Valid {
		item.Description = &description.String
	}
	if categoryID.Valid {
		item.CategoryID = &categoryID.String
	}

	h.logger.WithField("id", item.ID).Info("Stock item created")
	return &item, nil
}

// Update updates an existing stock item
func (h *DBHandler) Update(id string, req *models.StockItemUpdateRequest) (*models.StockItem, error) {
	// Validate unit if provided
	if req.Unit != nil && !models.IsValidUnit(*req.Unit) {
		return nil, fmt.Errorf("invalid unit: %s. Valid units are: %v", *req.Unit, models.ValidUnits)
	}

	query, err := h.queries.Get(stockItemSQL.UpdateStockItemQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var item models.StockItem
	var description, categoryID sql.NullString

	err = h.db.QueryRow(query, id, req.Name, req.Unit, req.Description, req.CategoryID).Scan(
		&item.ID, &item.Name, &item.Unit, &description, &categoryID, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update stock item: %w", err)
	}

	if description.Valid {
		item.Description = &description.String
	}
	if categoryID.Valid {
		item.CategoryID = &categoryID.String
	}

	h.logger.WithField("id", item.ID).Info("Stock item updated")
	return &item, nil
}

// Delete deletes a stock item
func (h *DBHandler) Delete(id string) error {
	checkQuery, err := h.queries.Get(stockItemSQL.CheckStockItemDependenciesQuery)
	if err != nil {
		return fmt.Errorf("failed to get check query: %w", err)
	}

	var count int
	if err := h.db.QueryRow(checkQuery, id).Scan(&count); err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete stock item: %d menu items use it as an ingredient", count)
	}

	deleteQuery, err := h.queries.Get(stockItemSQL.DeleteStockItemQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock item: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock item not found")
	}

	h.logger.WithField("id", id).Info("Stock item deleted")
	return nil
}
