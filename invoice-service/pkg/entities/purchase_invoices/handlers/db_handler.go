package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"invoice-service/pkg/entities/purchase_invoices/models"
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

func (h *DBHandler) Create(req *models.PurchaseInvoiceCreateRequest) (*models.PurchaseInvoice, error) {
	query := `
		INSERT INTO purchase_invoices (
			invoice_number, supplier_name, invoice_date, due_date,
			total_amount, status, image_url, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, invoice_number, supplier_name, invoice_date, due_date,
		          total_amount, status, image_url, notes, created_at, updated_at`

	var invoice models.PurchaseInvoice
	err := h.db.QueryRow(query,
		req.InvoiceNumber, req.SupplierName, req.InvoiceDate, req.DueDate,
		req.TotalAmount, req.Status, req.ImageURL, req.Notes,
	).Scan(
		&invoice.ID, &invoice.InvoiceNumber, &invoice.SupplierName, &invoice.InvoiceDate,
		&invoice.DueDate, &invoice.TotalAmount, &invoice.Status, &invoice.ImageURL,
		&invoice.Notes, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to create purchase invoice")
		return nil, fmt.Errorf("failed to create purchase invoice: %w", err)
	}

	h.logger.WithField("id", invoice.ID).Info("Purchase invoice created")
	return &invoice, nil
}

func (h *DBHandler) GetByID(id string) (*models.PurchaseInvoice, error) {
	query := `
		SELECT id, invoice_number, supplier_name, invoice_date, due_date,
		       total_amount, status, image_url, notes, created_at, updated_at
		FROM purchase_invoices WHERE id = $1`

	var invoice models.PurchaseInvoice
	err := h.db.QueryRow(query, id).Scan(
		&invoice.ID, &invoice.InvoiceNumber, &invoice.SupplierName, &invoice.InvoiceDate,
		&invoice.DueDate, &invoice.TotalAmount, &invoice.Status, &invoice.ImageURL,
		&invoice.Notes, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase invoice not found")
		}
		h.logger.WithError(err).Error("Failed to get purchase invoice")
		return nil, fmt.Errorf("failed to get purchase invoice: %w", err)
	}

	return &invoice, nil
}

func (h *DBHandler) Update(id string, req *models.PurchaseInvoiceUpdateRequest) (*models.PurchaseInvoice, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.SupplierName != nil {
		setParts = append(setParts, fmt.Sprintf("supplier_name = $%d", argIndex))
		args = append(args, *req.SupplierName)
		argIndex++
	}
	if req.InvoiceDate != nil {
		setParts = append(setParts, fmt.Sprintf("invoice_date = $%d", argIndex))
		args = append(args, *req.InvoiceDate)
		argIndex++
	}
	if req.DueDate != nil {
		setParts = append(setParts, fmt.Sprintf("due_date = $%d", argIndex))
		args = append(args, *req.DueDate)
		argIndex++
	}
	if req.TotalAmount != nil {
		setParts = append(setParts, fmt.Sprintf("total_amount = $%d", argIndex))
		args = append(args, *req.TotalAmount)
		argIndex++
	}
	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}
	if req.ImageURL != nil {
		setParts = append(setParts, fmt.Sprintf("image_url = $%d", argIndex))
		args = append(args, *req.ImageURL)
		argIndex++
	}
	if req.Notes != nil {
		setParts = append(setParts, fmt.Sprintf("notes = $%d", argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}

	if len(setParts) == 0 {
		return h.GetByID(id) // No updates, just return current
	}

	setParts = append(setParts, "updated_at = NOW()")
	query := fmt.Sprintf("UPDATE purchase_invoices SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := h.db.Exec(query, args...)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update purchase invoice")
		return nil, fmt.Errorf("failed to update purchase invoice: %w", err)
	}

	h.logger.WithField("id", id).Info("Purchase invoice updated")
	return h.GetByID(id)
}

func (h *DBHandler) Delete(id string) error {
	query := "DELETE FROM purchase_invoices WHERE id = $1"

	result, err := h.db.Exec(query, id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete purchase invoice")
		return fmt.Errorf("failed to delete purchase invoice: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase invoice not found")
	}

	h.logger.WithField("id", id).Info("Purchase invoice deleted")
	return nil
}

func (h *DBHandler) List(req *models.PurchaseInvoiceListRequest) (*models.PurchaseInvoiceListResponse, error) {
	whereParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.SupplierName != nil && *req.SupplierName != "" {
		whereParts = append(whereParts, fmt.Sprintf("supplier_name ILIKE $%d", argIndex))
		args = append(args, "%"+*req.SupplierName+"%")
		argIndex++
	}
	if req.Status != nil && *req.Status != "" {
		whereParts = append(whereParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM purchase_invoices %s", whereClause)
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		h.logger.WithError(err).Error("Failed to count purchase invoices")
		return nil, fmt.Errorf("failed to count purchase invoices: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, invoice_number, supplier_name, invoice_date, due_date,
		       total_amount, status, image_url, notes, created_at, updated_at
		FROM purchase_invoices %s
		ORDER BY invoice_date DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list purchase invoices")
		return nil, fmt.Errorf("failed to list purchase invoices: %w", err)
	}
	defer rows.Close()

	var invoices []models.PurchaseInvoice
	for rows.Next() {
		var invoice models.PurchaseInvoice
		err := rows.Scan(
			&invoice.ID, &invoice.InvoiceNumber, &invoice.SupplierName, &invoice.InvoiceDate,
			&invoice.DueDate, &invoice.TotalAmount, &invoice.Status, &invoice.ImageURL,
			&invoice.Notes, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			h.logger.WithError(err).Error("Failed to scan purchase invoice")
			return nil, fmt.Errorf("failed to scan purchase invoice: %w", err)
		}
		invoices = append(invoices, invoice)
	}

	return &models.PurchaseInvoiceListResponse{
		Invoices: invoices,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}