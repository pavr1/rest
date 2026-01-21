package handlers

import (
	"database/sql"
	"fmt"
	"strconv"

	"inventory-service/pkg/entities/stock_count/models"
	stockCountSQL "inventory-service/pkg/entities/stock_count/sql"
	sharedConfig "shared/config"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// parseFloat converts a string to float64, returns 0 on error
func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// DBHandler handles database operations for stock count
type DBHandler struct {
	db           *sharedDb.DbHandler
	queries      *stockCountSQL.Queries
	logger       *logrus.Logger
	config       *sharedConfig.Config
	portionGrams float64
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, config *sharedConfig.Config, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := stockCountSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	// Get default portion grams from config (default 120g)
	portionGrams := 120.0
	if config != nil {
		portionGrams = config.GetFloat("DEFAULT_PORTION_GRAMS")
		if portionGrams <= 0 {
			portionGrams = 120.0
		}
	}

	return &DBHandler{
		db:           db,
		queries:      queries,
		logger:       logger,
		config:       config,
		portionGrams: portionGrams,
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
	var invoiceID, unitPrice, costPerPortion sql.NullString
	var invoiceNumber, supplierName sql.NullString

	err = h.db.QueryRow(query, id).Scan(
		&sc.ID, &sc.StockVariantID, &invoiceID, &sc.Count, &sc.Unit,
		&unitPrice, &costPerPortion,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
		&sc.StockVariantName, &invoiceNumber, &supplierName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stock count record: %w", err)
	}

	// Handle nullable fields
	if invoiceID.Valid {
		sc.InvoiceID = &invoiceID.String
	}
	if unitPrice.Valid {
		val := parseFloat(unitPrice.String)
		sc.UnitPrice = &val
	}
	if costPerPortion.Valid {
		val := parseFloat(costPerPortion.String)
		sc.CostPerPortion = &val
	}
	if invoiceNumber.Valid {
		sc.InvoiceNumber = &invoiceNumber.String
	}
	if supplierName.Valid {
		sc.SupplierName = &supplierName.String
	}

	return &sc, nil
}

// Create creates a new stock count record
func (h *DBHandler) Create(req *models.StockCountCreateRequest) (*models.StockCount, error) {
	query, err := h.queries.Get(stockCountSQL.CreateStockCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	// Calculate cost per portion if unit_price is provided
	var costPerPortion *float64
	if req.UnitPrice != nil && *req.UnitPrice > 0 {
		totalKG, err := models.ConvertToKG(req.Count, req.Unit)
		if err != nil {
			return nil, fmt.Errorf("failed to convert unit: %w", err)
		}
		cost := models.CalculateCostPerPortion(totalKG, *req.UnitPrice, h.portionGrams)
		costPerPortion = &cost
	}

	var sc models.StockCount
	var invoiceID, unitPriceStr, costPerPortionStr sql.NullString

	err = h.db.QueryRow(query, req.StockVariantID, req.InvoiceID, req.Count, req.Unit, req.UnitPrice, costPerPortion, req.PurchasedAt).Scan(
		&sc.ID, &sc.StockVariantID, &invoiceID, &sc.Count, &sc.Unit,
		&unitPriceStr, &costPerPortionStr,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock count record: %w", err)
	}

	// Handle nullable fields
	if invoiceID.Valid {
		sc.InvoiceID = &invoiceID.String
	}
	if unitPriceStr.Valid {
		val := parseFloat(unitPriceStr.String)
		sc.UnitPrice = &val
	}
	if costPerPortionStr.Valid {
		val := parseFloat(costPerPortionStr.String)
		sc.CostPerPortion = &val
	}

	// Update avg_cost for the stock variant
	if err := h.UpdateAvgCost(req.StockVariantID); err != nil {
		h.logger.WithError(err).Warn("Failed to update avg_cost for stock variant")
	}

	h.logger.WithField("id", sc.ID).Info("Stock count record created")
	return &sc, nil
}

// Update updates an existing stock count record
func (h *DBHandler) Update(id string, req *models.StockCountUpdateRequest) (*models.StockCount, error) {
	// First get the existing record to have all values for cost calculation
	existing, err := h.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing record: %w", err)
	}
	if existing == nil {
		return nil, nil
	}

	// Determine the new values
	newCount := existing.Count
	newUnit := existing.Unit
	newUnitPrice := existing.UnitPrice

	if req.Count != nil {
		newCount = *req.Count
	}
	if req.Unit != nil {
		newUnit = *req.Unit
	}
	if req.UnitPrice != nil {
		newUnitPrice = req.UnitPrice
	}

	// Recalculate cost per portion if we have a unit price
	var costPerPortion *float64
	if newUnitPrice != nil && *newUnitPrice > 0 {
		totalKG, err := models.ConvertToKG(newCount, newUnit)
		if err != nil {
			return nil, fmt.Errorf("failed to convert unit: %w", err)
		}
		cost := models.CalculateCostPerPortion(totalKG, *newUnitPrice, h.portionGrams)
		costPerPortion = &cost
	}

	query, err := h.queries.Get(stockCountSQL.UpdateStockCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var sc models.StockCount
	var invoiceID, unitPriceStr, costPerPortionStr sql.NullString

	err = h.db.QueryRow(query, id, req.Count, req.Unit, newUnitPrice, costPerPortion, req.IsOut).Scan(
		&sc.ID, &sc.StockVariantID, &invoiceID, &sc.Count, &sc.Unit,
		&unitPriceStr, &costPerPortionStr,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update stock count record: %w", err)
	}

	// Handle nullable fields
	if invoiceID.Valid {
		sc.InvoiceID = &invoiceID.String
	}
	if unitPriceStr.Valid {
		val := parseFloat(unitPriceStr.String)
		sc.UnitPrice = &val
	}
	if costPerPortionStr.Valid {
		val := parseFloat(costPerPortionStr.String)
		sc.CostPerPortion = &val
	}

	// Update avg_cost for the stock variant
	if err := h.UpdateAvgCost(sc.StockVariantID); err != nil {
		h.logger.WithError(err).Warn("Failed to update avg_cost for stock variant")
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
	var invoiceID, unitPriceStr, costPerPortionStr sql.NullString

	err = h.db.QueryRow(query, id, isOut).Scan(
		&sc.ID, &sc.StockVariantID, &invoiceID, &sc.Count, &sc.Unit,
		&unitPriceStr, &costPerPortionStr,
		&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to mark stock out: %w", err)
	}

	// Handle nullable fields
	if invoiceID.Valid {
		sc.InvoiceID = &invoiceID.String
	}
	if unitPriceStr.Valid {
		val := parseFloat(unitPriceStr.String)
		sc.UnitPrice = &val
	}
	if costPerPortionStr.Valid {
		val := parseFloat(costPerPortionStr.String)
		sc.CostPerPortion = &val
	}

	// Update avg_cost for the stock variant (since is_out affects the avg calculation)
	if err := h.UpdateAvgCost(sc.StockVariantID); err != nil {
		h.logger.WithError(err).Warn("Failed to update avg_cost for stock variant")
	}

	h.logger.WithFields(logrus.Fields{"id": sc.ID, "is_out": isOut}).Info("Stock count out status updated")
	return &sc, nil
}

// Delete deletes a stock count record
func (h *DBHandler) Delete(id string) error {
	// First get the stock variant ID for updating avg_cost after deletion
	existing, err := h.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get existing record: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("stock count record not found")
	}
	stockVariantID := existing.StockVariantID

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

	// Update avg_cost for the stock variant after deletion
	if err := h.UpdateAvgCost(stockVariantID); err != nil {
		h.logger.WithError(err).Warn("Failed to update avg_cost for stock variant")
	}

	h.logger.WithField("id", id).Info("Stock count record deleted")
	return nil
}

// UpdateAvgCost updates the average cost per portion for a stock variant
func (h *DBHandler) UpdateAvgCost(stockVariantID string) error {
	query, err := h.queries.Get(stockCountSQL.CalculateAvgCostQuery)
	if err != nil {
		return fmt.Errorf("failed to get calculate_avg_cost query: %w", err)
	}

	var id string
	var avgCost float64
	err = h.db.QueryRow(query, stockVariantID).Scan(&id, &avgCost)
	if err != nil {
		return fmt.Errorf("failed to update avg_cost: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"stock_variant_id": stockVariantID,
		"avg_cost":         avgCost,
	}).Info("Stock variant avg_cost updated")
	return nil
}

// scanStockCounts scans multiple stock count rows
func (h *DBHandler) scanStockCounts(rows *sql.Rows) ([]models.StockCount, error) {
	var stockCounts []models.StockCount
	for rows.Next() {
		var sc models.StockCount
		var invoiceID, unitPriceStr, costPerPortionStr sql.NullString
		var invoiceNumber, supplierName sql.NullString

		if err := rows.Scan(
			&sc.ID, &sc.StockVariantID, &invoiceID, &sc.Count, &sc.Unit,
			&unitPriceStr, &costPerPortionStr,
			&sc.PurchasedAt, &sc.IsOut, &sc.CreatedAt, &sc.UpdatedAt,
			&sc.StockVariantName, &invoiceNumber, &supplierName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan stock count record: %w", err)
		}

		// Handle nullable fields
		if invoiceID.Valid {
			sc.InvoiceID = &invoiceID.String
		}
		if unitPriceStr.Valid {
			val := parseFloat(unitPriceStr.String)
			sc.UnitPrice = &val
		}
		if costPerPortionStr.Valid {
			val := parseFloat(costPerPortionStr.String)
			sc.CostPerPortion = &val
		}
		if invoiceNumber.Valid {
			sc.InvoiceNumber = &invoiceNumber.String
		}
		if supplierName.Valid {
			sc.SupplierName = &supplierName.String
		}

		stockCounts = append(stockCounts, sc)
	}
	return stockCounts, nil
}
