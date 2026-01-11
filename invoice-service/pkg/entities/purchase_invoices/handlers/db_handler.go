package handlers

import (
	"database/sql"
	"fmt"

	"invoice-service/pkg/entities/purchase_invoices/models"
	"invoice-service/pkg/entities/purchase_invoices/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

type DBHandler struct {
	db      *sharedDb.DbHandler
	logger  *logrus.Logger
	queries *sql.Queries
}

func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := sql.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		logger:  logger,
		queries: queries,
	}, nil
}

func (h *DBHandler) Create(req *models.PurchaseInvoiceCreateRequest) (*models.PurchaseInvoice, error) {
	query, err := h.queries.Get(sql.CreatePurchaseInvoiceQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get create query: %w", err)
	}

	var invoice models.PurchaseInvoice
	err = h.db.QueryRow(query,
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
	query, err := h.queries.Get(sql.GetPurchaseInvoiceQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var invoice models.PurchaseInvoice
	err = h.db.QueryRow(query, id).Scan(
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
	query, err := h.queries.Get(sql.UpdatePurchaseInvoiceQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get update query: %w", err)
	}

	_, err = h.db.Exec(query,
		req.SupplierName, req.InvoiceDate, req.DueDate, req.TotalAmount,
		req.Status, req.ImageURL, req.Notes, id,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update purchase invoice")
		return nil, fmt.Errorf("failed to update purchase invoice: %w", err)
	}

	h.logger.WithField("id", id).Info("Purchase invoice updated")
	return h.GetByID(id)
}

func (h *DBHandler) Delete(id string) error {
	query, err := h.queries.Get(sql.DeletePurchaseInvoiceQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

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
	// Get the list and count queries
	listQuery, err := h.queries.Get(sql.ListPurchaseInvoicesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	countQuery, err := h.queries.Get(sql.CountPurchaseInvoicesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	// Get total count
	var total int
	err = h.db.QueryRow(countQuery, req.SupplierName, req.Status).Scan(&total)
	if err != nil {
		h.logger.WithError(err).Error("Failed to count purchase invoices")
		return nil, fmt.Errorf("failed to count purchase invoices: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	rows, err := h.db.Query(listQuery, req.SupplierName, req.Status, req.Limit, offset)
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
