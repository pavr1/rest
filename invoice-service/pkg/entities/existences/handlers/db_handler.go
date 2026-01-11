package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	"invoice-service/pkg/entities/existences/models"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

type DBHandler struct {
	db     *sharedDb.DbHandler
	logger *logrus.Logger
}

func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	return &DBHandler{
		db:     db,
		logger: logger,
	}, nil
}

func (h *DBHandler) Create(req *models.ExistenceCreateRequest) (*models.Existence, error) {
	// First get the invoice detail to get the necessary information
	invoiceDetailQuery := `
		SELECT id.stock_item_id, si.name, id.quantity, id.unit_price,
		       id.unit_of_measure, id.items_per_unit, id.expiry_date, id.batch_number
		FROM invoice_details id
		LEFT JOIN stock_items si ON id.stock_item_id = si.id
		WHERE id.id = $1`

	var stockItemID *string
	var stockItemName *string
	var quantity, unitPrice float64
	var unitOfMeasure string
	var itemsPerUnit *float64
	var expiryDate, batchNumber *string

	err := h.db.QueryRow(invoiceDetailQuery, req.InvoiceDetailID).Scan(
		&stockItemID, &stockItemName, &quantity, &unitPrice,
		&unitOfMeasure, &itemsPerUnit, &expiryDate, &batchNumber,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get invoice detail for existence creation")
		return nil, fmt.Errorf("failed to get invoice detail: %w", err)
	}

	if stockItemID == nil {
		return nil, fmt.Errorf("invoice detail must have a stock item assigned")
	}

	// Calculate units purchased and cost per unit
	var unitsPurchased float64
	var costPerUnit float64

	if itemsPerUnit != nil && *itemsPerUnit > 0 {
		// If items_per_unit is specified (e.g., 12 breadsticks per unit)
		unitsPurchased = quantity * *itemsPerUnit
		costPerUnit = unitPrice / *itemsPerUnit
	} else {
		// Direct unit measurement
		unitsPurchased = quantity
		costPerUnit = unitPrice
	}

	totalCost := unitsPurchased * costPerUnit

	// Use request expiry/batch if provided, otherwise use invoice detail values
	finalExpiryDate := req.ExpiryDate
	finalBatchNumber := req.BatchNumber

	if finalExpiryDate == nil && expiryDate != nil {
		// Convert string to time.Time if needed (would need proper parsing)
		finalExpiryDate = req.ExpiryDate // Keep request value
	}
	if finalBatchNumber == nil && batchNumber != nil {
		finalBatchNumber = req.BatchNumber // Keep request value
	}

	query := `
		INSERT INTO existences (
			invoice_detail_id, stock_item_id, units_purchased, cost_per_unit,
			total_cost, expiry_date, batch_number, current_stock
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, invoice_detail_id, stock_item_id, units_purchased,
		          cost_per_unit, total_cost, expiry_date, batch_number,
		          current_stock, created_at, updated_at`

	var existence models.Existence
	err = h.db.QueryRow(query,
		req.InvoiceDetailID, stockItemID, unitsPurchased, costPerUnit,
		totalCost, finalExpiryDate, finalBatchNumber, unitsPurchased, // current_stock starts as units_purchased
	).Scan(
		&existence.ID, &existence.InvoiceDetailID, &existence.StockItemID,
		&existence.UnitsPurchased, &existence.CostPerUnit, &existence.TotalCost,
		&existence.ExpiryDate, &existence.BatchNumber, &existence.CurrentStock,
		&existence.CreatedAt, &existence.UpdatedAt,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to create existence")
		return nil, fmt.Errorf("failed to create existence: %w", err)
	}

	// Update the stock item's current stock and unit cost
	err = h.updateStockItemFromExistence(*stockItemID, costPerUnit, unitsPurchased)
	if err != nil {
		h.logger.WithError(err).Warn("Failed to update stock item, but existence was created")
	}

	existence.StockItemName = stockItemName

	h.logger.WithField("id", existence.ID).Info("Existence created and stock item updated")
	return &existence, nil
}

// updateStockItemFromExistence updates the stock item's unit cost and current stock
func (h *DBHandler) updateStockItemFromExistence(stockItemID string, newCostPerUnit, unitsPurchased float64) error {
	// Get current stock item info
	query := `SELECT current_stock, unit_cost FROM stock_items WHERE id = $1`
	var currentStock float64
	var currentUnitCost *float64

	err := h.db.QueryRow(query, stockItemID).Scan(&currentStock, &currentUnitCost)
	if err != nil {
		return fmt.Errorf("failed to get current stock item: %w", err)
	}

	// Calculate new values
	newTotalStock := currentStock + unitsPurchased
	var newAverageCost *float64

	if currentUnitCost != nil {
		// Calculate weighted average cost
		currentValue := currentStock * *currentUnitCost
		newValue := unitsPurchased * newCostPerUnit
		totalValue := currentValue + newValue
		averageCost := totalValue / newTotalStock
		newAverageCost = &averageCost
	} else {
		// First time setting cost
		newAverageCost = &newCostPerUnit
	}

	// Update stock item
	updateQuery := `
		UPDATE stock_items
		SET current_stock = $1, unit_cost = $2, total_value = $3, updated_at = NOW()
		WHERE id = $4`

	totalValue := newTotalStock * *newAverageCost
	_, err = h.db.Exec(updateQuery, newTotalStock, newAverageCost, totalValue, stockItemID)
	if err != nil {
		return fmt.Errorf("failed to update stock item: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"stock_item_id": stockItemID,
		"new_stock":     newTotalStock,
		"unit_cost":     *newAverageCost,
	}).Info("Stock item updated from existence")

	return nil
}

func (h *DBHandler) GetByID(id string) (*models.Existence, error) {
	query := `
		SELECT e.id, e.invoice_detail_id, e.stock_item_id, si.name as stock_item_name,
		       e.units_purchased, e.cost_per_unit, e.total_cost, e.expiry_date,
		       e.batch_number, e.current_stock, e.created_at, e.updated_at
		FROM existences e
		LEFT JOIN stock_items si ON e.stock_item_id = si.id
		WHERE e.id = $1`

	var existence models.Existence
	err := h.db.QueryRow(query, id).Scan(
		&existence.ID, &existence.InvoiceDetailID, &existence.StockItemID,
		&existence.StockItemName, &existence.UnitsPurchased, &existence.CostPerUnit,
		&existence.TotalCost, &existence.ExpiryDate, &existence.BatchNumber,
		&existence.CurrentStock, &existence.CreatedAt, &existence.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("existence not found")
		}
		h.logger.WithError(err).Error("Failed to get existence")
		return nil, fmt.Errorf("failed to get existence: %w", err)
	}

	return &existence, nil
}

func (h *DBHandler) Update(id string, req *models.ExistenceUpdateRequest) (*models.Existence, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.CurrentStock != nil {
		setParts = append(setParts, fmt.Sprintf("current_stock = $%d", argIndex))
		args = append(args, *req.CurrentStock)
		argIndex++
	}
	if req.ExpiryDate != nil {
		setParts = append(setParts, fmt.Sprintf("expiry_date = $%d", argIndex))
		args = append(args, *req.ExpiryDate)
		argIndex++
	}
	if req.BatchNumber != nil {
		setParts = append(setParts, fmt.Sprintf("batch_number = $%d", argIndex))
		args = append(args, *req.BatchNumber)
		argIndex++
	}

	if len(setParts) == 0 {
		return h.GetByID(id) // No updates, just return current
	}

	setParts = append(setParts, "updated_at = NOW()")
	query := fmt.Sprintf("UPDATE existences SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := h.db.Exec(query, args...)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update existence")
		return nil, fmt.Errorf("failed to update existence: %w", err)
	}

	h.logger.WithField("id", id).Info("Existence updated")
	return h.GetByID(id)
}

func (h *DBHandler) Delete(id string) error {
	query := "DELETE FROM existences WHERE id = $1"

	result, err := h.db.Exec(query, id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete existence")
		return fmt.Errorf("failed to delete existence: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("existence not found")
	}

	h.logger.WithField("id", id).Info("Existence deleted")
	return nil
}

func (h *DBHandler) List(req *models.ExistenceListRequest) (*models.ExistenceListResponse, error) {
	whereParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.StockItemID != nil && *req.StockItemID != "" {
		whereParts = append(whereParts, fmt.Sprintf("e.stock_item_id = $%d", argIndex))
		args = append(args, *req.StockItemID)
		argIndex++
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM existences e %s", whereClause)
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		h.logger.WithError(err).Error("Failed to count existences")
		return nil, fmt.Errorf("failed to count existences: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT e.id, e.invoice_detail_id, e.stock_item_id, si.name as stock_item_name,
		       e.units_purchased, e.cost_per_unit, e.total_cost, e.expiry_date,
		       e.batch_number, e.current_stock, e.created_at, e.updated_at
		FROM existences e
		LEFT JOIN stock_items si ON e.stock_item_id = si.id
		%s
		ORDER BY e.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list existences")
		return nil, fmt.Errorf("failed to list existences: %w", err)
	}
	defer rows.Close()

	var existences []models.Existence
	for rows.Next() {
		var existence models.Existence
		err := rows.Scan(
			&existence.ID, &existence.InvoiceDetailID, &existence.StockItemID,
			&existence.StockItemName, &existence.UnitsPurchased, &existence.CostPerUnit,
			&existence.TotalCost, &existence.ExpiryDate, &existence.BatchNumber,
			&existence.CurrentStock, &existence.CreatedAt, &existence.UpdatedAt,
		)
		if err != nil {
			h.logger.WithError(err).Error("Failed to scan existence")
			return nil, fmt.Errorf("failed to scan existence: %w", err)
		}
		existences = append(existences, existence)
	}

	return &models.ExistenceListResponse{
		Existences: existences,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}