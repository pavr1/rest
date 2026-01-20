package handlers

import (
	"database/sql"
	"fmt"
	"inventory-service/pkg/entities/stock_count/models"
	stockCountSQL "inventory-service/pkg/entities/stock_count/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for stock count
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *stockCountSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := stockCountSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of all stock count records
func (h *DBHandler) List(page, limit int) (*models.StockCountListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockCountSQL.CountStockCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock count records: %w", err)
	}

	listQuery, err := h.queries.Get(stockCountSQL.ListStockCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock count records: %w", err)
	}
	defer rows.Close()

	stockCounts, err := h.scanStockCounts(rows)
	if err != nil {
		return nil, err
	}

	return &models.StockCountListResponse{
		StockCounts: stockCounts,
		Total:       total,
		Page:        page,
		Limit:       limit,
	}, nil
}

// ListByVariant returns stock count records for a specific variant
func (h *DBHandler) ListByVariant(variantID string, page, limit int) (*models.StockCountListResponse, error) {
	offset := (page - 1) * limit

	countQuery, err := h.queries.Get(stockCountSQL.CountStockCountByVariantQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, variantID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count stock count records: %w", err)
	}

	listQuery, err := h.queries.Get(stockCountSQL.ListStockCountByVariantQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, variantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock count records: %w", err)
	}
	defer rows.Close()

	stockCounts, err := h.scanStockCounts(rows)
	if err != nil {
		return nil, err
	}

	return &models.StockCountListResponse{
		StockCounts: stockCounts,
		Total:       total,
		Page:        page,
		Limit:       limit,
	}, nil
}

// GetByID returns a stock count record by ID
func (h *DBHandler) GetByID(id string) (*models.StockCount, error) {
	query, err := h.queries.Get(stockCountSQL.GetStockCountByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var sc models.StockCount
	err = h.db.QueryRow(query, id).Scan(
		&sc.ID, &sc.StockVariantID, &sc.InvoiceID, &sc.Count, &sc.Unit,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
		&sc.StockVariantName, &sc.InvoiceNumber, &sc.SupplierName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock count record: %w", err)
	}

	return &sc, nil
}

// Create creates a new stock count record
func (h *DBHandler) Create(req *models.StockCountCreateRequest) (*models.StockCount, error) {
	query, err := h.queries.Get(stockCountSQL.CreateStockCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var sc models.StockCount
	err = h.db.QueryRow(query, req.StockVariantID, req.InvoiceID, req.Count, req.Unit, req.PurchasedAt).Scan(
		&sc.ID, &sc.StockVariantID, &sc.InvoiceID, &sc.Count, &sc.Unit,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock count record: %w", err)
	}

	h.logger.WithField("id", sc.ID).Info("Stock count record created")
	return &sc, nil
}

// Update updates an existing stock count record
func (h *DBHandler) Update(id string, req *models.StockCountUpdateRequest) (*models.StockCount, error) {
	query, err := h.queries.Get(stockCountSQL.UpdateStockCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var sc models.StockCount
	err = h.db.QueryRow(query, id, req.Count, req.Unit, req.IsOut).Scan(
		&sc.ID, &sc.StockVariantID, &sc.InvoiceID, &sc.Count, &sc.Unit,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update stock count record: %w", err)
	}

	h.logger.WithField("id", sc.ID).Info("Stock count record updated")
	return &sc, nil
}

// MarkOut marks a stock count record as out/available
func (h *DBHandler) MarkOut(id string, isOut bool) (*models.StockCount, error) {
	query, err := h.queries.Get(stockCountSQL.MarkStockOutQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var sc models.StockCount
	err = h.db.QueryRow(query, id, isOut).Scan(
		&sc.ID, &sc.StockVariantID, &sc.InvoiceID, &sc.Count, &sc.Unit,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to mark stock out: %w", err)
	}

	h.logger.WithFields(logrus.Fields{"id": sc.ID, "is_out": isOut}).Info("Stock count out status updated")
	return &sc, nil
}

// Delete deletes a stock count record
func (h *DBHandler) Delete(id string) error {
	query, err := h.queries.Get(stockCountSQL.DeleteStockCountQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock count record: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("stock count record not found")
	}

	h.logger.WithField("id", id).Info("Stock count record deleted")
	return nil
}

// scanStockCounts scans multiple stock count rows
func (h *DBHandler) scanStockCounts(rows *sql.Rows) ([]models.StockCount, error) {
	var stockCounts []models.StockCount
	for rows.Next() {
		var sc models.StockCount
		if err := rows.Scan(
			&sc.ID, &sc.StockVariantID, &sc.InvoiceID, &sc.Count, &sc.Unit,
			&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
			&sc.StockVariantName, &sc.InvoiceNumber, &sc.SupplierName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan stock count record: %w", err)
		}
		stockCounts = append(stockCounts, sc)
	}
	return stockCounts, nil
}
