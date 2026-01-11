package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	"invoice-service/pkg/entities/invoice_details/models"
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

func (h *DBHandler) Create(invoiceID string, req *models.InvoiceDetailCreateRequest) (*models.InvoiceDetail, error) {
	// Calculate total price
	totalPrice := req.Quantity * req.UnitPrice

	query := `
		INSERT INTO invoice_details (
			invoice_id, stock_item_id, description, quantity,
			unit_of_measure, items_per_unit, unit_price, total_price,
			expiry_date, batch_number
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, invoice_id, stock_item_id, description, quantity,
		          unit_of_measure, items_per_unit, unit_price, total_price,
		          expiry_date, batch_number, created_at, updated_at`

	var detail models.InvoiceDetail
	err := h.db.QueryRow(query,
		invoiceID, req.StockItemID, req.Description, req.Quantity,
		req.UnitOfMeasure, req.ItemsPerUnit, req.UnitPrice, totalPrice,
		req.ExpiryDate, req.BatchNumber,
	).Scan(
		&detail.ID, &detail.InvoiceID, &detail.StockItemID, &detail.Description,
		&detail.Quantity, &detail.UnitOfMeasure, &detail.ItemsPerUnit,
		&detail.UnitPrice, &detail.TotalPrice, &detail.ExpiryDate,
		&detail.BatchNumber, &detail.CreatedAt, &detail.UpdatedAt,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to create invoice detail")
		return nil, fmt.Errorf("failed to create invoice detail: %w", err)
	}

	h.logger.WithField("id", detail.ID).Info("Invoice detail created")
	return &detail, nil
}

func (h *DBHandler) GetByID(invoiceID, id string) (*models.InvoiceDetail, error) {
	query := `
		SELECT id, invoice_id, stock_item_id, description, quantity,
		       unit_of_measure, items_per_unit, unit_price, total_price,
		       expiry_date, batch_number, created_at, updated_at
		FROM invoice_details WHERE id = $1 AND invoice_id = $2`

	var detail models.InvoiceDetail
	err := h.db.QueryRow(query, id, invoiceID).Scan(
		&detail.ID, &detail.InvoiceID, &detail.StockItemID, &detail.Description,
		&detail.Quantity, &detail.UnitOfMeasure, &detail.ItemsPerUnit,
		&detail.UnitPrice, &detail.TotalPrice, &detail.ExpiryDate,
		&detail.BatchNumber, &detail.CreatedAt, &detail.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice detail not found")
		}
		h.logger.WithError(err).Error("Failed to get invoice detail")
		return nil, fmt.Errorf("failed to get invoice detail: %w", err)
	}

	return &detail, nil
}

func (h *DBHandler) Update(invoiceID, id string, req *models.InvoiceDetailUpdateRequest) (*models.InvoiceDetail, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.StockItemID != nil {
		setParts = append(setParts, fmt.Sprintf("stock_item_id = $%d", argIndex))
		args = append(args, *req.StockItemID)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Quantity != nil {
		setParts = append(setParts, fmt.Sprintf("quantity = $%d", argIndex))
		args = append(args, *req.Quantity)
		argIndex++
	}
	if req.UnitOfMeasure != nil {
		setParts = append(setParts, fmt.Sprintf("unit_of_measure = $%d", argIndex))
		args = append(args, *req.UnitOfMeasure)
		argIndex++
	}
	if req.ItemsPerUnit != nil {
		setParts = append(setParts, fmt.Sprintf("items_per_unit = $%d", argIndex))
		args = append(args, *req.ItemsPerUnit)
		argIndex++
	}
	if req.UnitPrice != nil {
		setParts = append(setParts, fmt.Sprintf("unit_price = $%d", argIndex))
		args = append(args, *req.UnitPrice)
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

	// Recalculate total price if quantity or unit price changed
	if req.Quantity != nil || req.UnitPrice != nil {
		// Get current values if not provided
		current, err := h.GetByID(invoiceID, id)
		if err != nil {
			return nil, err
		}

		quantity := current.Quantity
		unitPrice := current.UnitPrice

		if req.Quantity != nil {
			quantity = *req.Quantity
		}
		if req.UnitPrice != nil {
			unitPrice = *req.UnitPrice
		}

		totalPrice := quantity * unitPrice
		setParts = append(setParts, fmt.Sprintf("total_price = $%d", argIndex))
		args = append(args, totalPrice)
		argIndex++
	}

	if len(setParts) == 0 {
		return h.GetByID(invoiceID, id) // No updates, just return current
	}

	setParts = append(setParts, "updated_at = NOW()")
	query := fmt.Sprintf("UPDATE invoice_details SET %s WHERE id = $%d AND invoice_id = $%d",
		strings.Join(setParts, ", "), argIndex, argIndex+1)
	args = append(args, id, invoiceID)

	_, err := h.db.Exec(query, args...)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update invoice detail")
		return nil, fmt.Errorf("failed to update invoice detail: %w", err)
	}

	h.logger.WithField("id", id).Info("Invoice detail updated")
	return h.GetByID(invoiceID, id)
}

func (h *DBHandler) Delete(invoiceID, id string) error {
	query := "DELETE FROM invoice_details WHERE id = $1 AND invoice_id = $2"

	result, err := h.db.Exec(query, id, invoiceID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete invoice detail")
		return fmt.Errorf("failed to delete invoice detail: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invoice detail not found")
	}

	h.logger.WithField("id", id).Info("Invoice detail deleted")
	return nil
}

func (h *DBHandler) ListByInvoice(invoiceID string) (*models.InvoiceDetailListResponse, error) {
	query := `
		SELECT id, invoice_id, stock_item_id, si.name as stock_item_name,
		       description, quantity, unit_of_measure, items_per_unit,
		       unit_price, total_price, expiry_date, batch_number,
		       id.created_at, id.updated_at
		FROM invoice_details id
		LEFT JOIN stock_items si ON id.stock_item_id = si.id
		WHERE id.invoice_id = $1
		ORDER BY id.created_at`

	rows, err := h.db.Query(query, invoiceID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list invoice details")
		return nil, fmt.Errorf("failed to list invoice details: %w", err)
	}
	defer rows.Close()

	var details []models.InvoiceDetail
	for rows.Next() {
		var detail models.InvoiceDetail
		err := rows.Scan(
			&detail.ID, &detail.InvoiceID, &detail.StockItemID, &detail.StockItemName,
			&detail.Description, &detail.Quantity, &detail.UnitOfMeasure, &detail.ItemsPerUnit,
			&detail.UnitPrice, &detail.TotalPrice, &detail.ExpiryDate,
			&detail.BatchNumber, &detail.CreatedAt, &detail.UpdatedAt,
		)
		if err != nil {
			h.logger.WithError(err).Error("Failed to scan invoice detail")
			return nil, fmt.Errorf("failed to scan invoice detail: %w", err)
		}
		details = append(details, detail)
	}

	return &models.InvoiceDetailListResponse{
		InvoiceID:       invoiceID,
		InvoiceDetails: details,
	}, nil
}