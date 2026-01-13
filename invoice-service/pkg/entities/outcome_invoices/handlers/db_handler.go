package handlers

import (
	"database/sql"
	"fmt"

	"invoice-service/pkg/entities/outcome_invoices/models"
	outcomesql "invoice-service/pkg/entities/outcome_invoices/sql"
	sharedDb "shared/db"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type DBHandler struct {
	db      *sharedDb.DbHandler
	logger  *logrus.Logger
	queries *outcomesql.Queries
}

func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := outcomesql.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		logger:  logger,
		queries: queries,
	}, nil
}

func (h *DBHandler) Create(req *models.OutcomeInvoiceCreateRequest) (*models.OutcomeInvoice, error) {
	query, err := h.queries.Get(outcomesql.CreateOutcomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get create query: %w", err)
	}

	var invoice models.OutcomeInvoice
	err = h.db.QueryRow(query,
		req.InvoiceNumber,
		req.SupplierID,
		req.TransactionDate,
		req.TotalAmount,
		req.ImageURL,
		req.Notes,
	).Scan(
		&invoice.ID, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to create outcome invoice")
		return nil, fmt.Errorf("failed to create outcome invoice: %w", err)
	}

	// Fill in the rest of the fields from the request
	invoice.InvoiceNumber = req.InvoiceNumber
	invoice.SupplierID = req.SupplierID
	invoice.TransactionDate = req.TransactionDate
	invoice.TotalAmount = req.TotalAmount
	invoice.ImageURL = req.ImageURL
	invoice.Notes = req.Notes

	h.logger.WithField("id", invoice.ID).Info("Outcome invoice created")
	return &invoice, nil
}

func (h *DBHandler) GetByID(id string) (*models.OutcomeInvoice, error) {
	query, err := h.queries.Get(outcomesql.GetOutcomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var invoice models.OutcomeInvoice
	err = h.db.QueryRow(query, id).Scan(
		&invoice.ID,
		&invoice.InvoiceNumber,
		&invoice.SupplierID,
		&invoice.TransactionDate,
		&invoice.TotalAmount,
		&invoice.ImageURL,
		&invoice.Notes,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("outcome invoice not found")
		}
		h.logger.WithError(err).Error("Failed to get outcome invoice")
		return nil, fmt.Errorf("failed to get outcome invoice: %w", err)
	}

	return &invoice, nil
}

func (h *DBHandler) Update(id string, req *models.OutcomeInvoiceUpdateRequest) (*models.OutcomeInvoice, error) {
	query, err := h.queries.Get(outcomesql.UpdateOutcomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get update query: %w", err)
	}

	_, err = h.db.Exec(query,
		req.SupplierID,
		req.TransactionDate,
		req.TotalAmount,
		req.ImageURL,
		req.Notes,
		id,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update outcome invoice")
		return nil, fmt.Errorf("failed to update outcome invoice: %w", err)
	}

	h.logger.WithField("id", id).Info("Outcome invoice updated")
	return h.GetByID(id)
}

func (h *DBHandler) Delete(id string) error {
	query, err := h.queries.Get(outcomesql.DeleteOutcomeInvoice)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(query, id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete outcome invoice")
		return fmt.Errorf("failed to delete outcome invoice: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("outcome invoice not found")
	}

	h.logger.WithField("id", id).Info("Outcome invoice deleted")
	return nil
}

func (h *DBHandler) List(req *models.OutcomeInvoiceListRequest) (*models.OutcomeInvoiceListResponse, error) {
	// Get the list and count queries
	listQuery, err := h.queries.Get(outcomesql.ListOutcomeInvoices)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	countQuery, err := h.queries.Get(outcomesql.CountOutcomeInvoices)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	// Get total count
	var total int
	err = h.db.QueryRow(countQuery, req.SupplierID).Scan(&total)
	if err != nil {
		h.logger.WithError(err).Error("Failed to count outcome invoices")
		return nil, fmt.Errorf("failed to count outcome invoices: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	rows, err := h.db.Query(listQuery, req.SupplierID, req.Limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list outcome invoices")
		return nil, fmt.Errorf("failed to list outcome invoices: %w", err)
	}
	defer rows.Close()

	var invoices []models.OutcomeInvoice
	for rows.Next() {
		var invoice models.OutcomeInvoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNumber,
			&invoice.SupplierID,
			&invoice.TransactionDate,
			&invoice.TotalAmount,
			&invoice.ImageURL,
			&invoice.Notes,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			h.logger.WithError(err).Error("Failed to scan outcome invoice")
			return nil, fmt.Errorf("failed to scan outcome invoice: %w", err)
		}
		invoices = append(invoices, invoice)
	}

	return &models.OutcomeInvoiceListResponse{
		Invoices: invoices,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}
