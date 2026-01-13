package handlers

import (
	"context"
	"database/sql"
	"fmt"

	invoiceItemSql "invoice-service/pkg/entities/invoice_items/sql"
	"invoice-service/pkg/entities/outcome_invoices/models"
	outcomesql "invoice-service/pkg/entities/outcome_invoices/sql"
	sharedDb "shared/db"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type DBHandler struct {
	db                 *sharedDb.DbHandler
	logger             *logrus.Logger
	queries            *outcomesql.Queries
	invoiceItemQueries *invoiceItemSql.Queries
}

func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := outcomesql.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load outcome invoice SQL queries: %w", err)
	}

	invoiceItemQueries, err := invoiceItemSql.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load invoice item SQL queries: %w", err)
	}

	return &DBHandler{
		db:                 db,
		logger:             logger,
		queries:            queries,
		invoiceItemQueries: invoiceItemQueries,
	}, nil
}

func (h *DBHandler) Create(req *models.OutcomeInvoiceCreateRequest) (*models.OutcomeInvoice, error) {
	// Start transaction
	tx, err := h.db.BeginTx(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Calculate total from items if not provided
	total := req.TotalAmount
	if total == nil && req.InvoiceItems != nil {
		calculatedTotal := 0.0
		for _, item := range req.InvoiceItems {
			calculatedTotal += item.Count * item.Price
		}
		total = &calculatedTotal
	}

	query, err := h.queries.Get(outcomesql.CreateOutcomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get create query: %w", err)
	}

	var invoice models.OutcomeInvoice
	err = tx.QueryRow(query,
		req.InvoiceNumber,
		req.SupplierID,
		req.TransactionDate,
		total,
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
	invoice.TotalAmount = total
	invoice.ImageURL = req.ImageURL
	invoice.Notes = req.Notes

	// Create invoice items if provided
	if req.InvoiceItems != nil && len(req.InvoiceItems) > 0 {
		for _, itemReq := range req.InvoiceItems {
			item, err := h.createInvoiceItem(tx, invoice.ID, "outcome", &itemReq)
			if err != nil {
				h.logger.WithError(err).Error("Failed to create invoice item")
				return nil, fmt.Errorf("failed to create invoice item: %w", err)
			}
			invoice.InvoiceItems = append(invoice.InvoiceItems, *item)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	h.logger.WithField("id", invoice.ID).Info("Outcome invoice created with items")
	return &invoice, nil
}

// createInvoiceItem creates an invoice item within a transaction
func (h *DBHandler) createInvoiceItem(tx *sql.Tx, invoiceID, invoiceType string, req *models.InvoiceItemCreateRequest) (*models.InvoiceItem, error) {
	// Get the create query from invoice items SQL
	query, err := h.invoiceItemQueries.Get(invoiceItemSql.CreateInvoiceItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice item create query: %w", err)
	}

	var item models.InvoiceItem
	err = tx.QueryRow(query,
		invoiceID, invoiceType, req.Detail, req.Count, req.UnitType,
		req.Price, req.ItemsPerUnit, req.ExpirationDate,
	).Scan(
		&item.ID, &item.Total, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create invoice item: %w", err)
	}

	// Fill in the rest of the fields
	item.InvoiceID = invoiceID
	item.InvoiceType = invoiceType
	item.Detail = req.Detail
	item.Count = req.Count
	item.UnitType = req.UnitType
	item.Price = req.Price
	item.ItemsPerUnit = req.ItemsPerUnit
	item.ExpirationDate = req.ExpirationDate

	return &item, nil
}

// getInvoiceItems retrieves all invoice items for a given invoice ID
func (h *DBHandler) getInvoiceItems(invoiceID string) ([]models.InvoiceItem, error) {
	query, err := h.invoiceItemQueries.Get(invoiceItemSql.ListInvoiceItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(query, invoiceID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to query invoice items")
		return nil, fmt.Errorf("failed to query invoice items: %w", err)
	}
	defer rows.Close()

	var items []models.InvoiceItem
	for rows.Next() {
		var item models.InvoiceItem
		err := rows.Scan(
			&item.ID,
			&item.InvoiceID,
			&item.InvoiceType,
			&item.Detail,
			&item.Count,
			&item.UnitType,
			&item.Price,
			&item.ItemsPerUnit,
			&item.Total,
			&item.ExpirationDate,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			h.logger.WithError(err).Error("Failed to scan invoice item")
			return nil, fmt.Errorf("failed to scan invoice item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invoice items: %w", err)
	}

	return items, nil
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

	// Fetch invoice items for this invoice
	invoiceItems, err := h.getInvoiceItems(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}
	invoice.InvoiceItems = invoiceItems

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
