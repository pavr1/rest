package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"invoice-service/pkg/entities/income_invoices/models"
	incomesql "invoice-service/pkg/entities/income_invoices/sql"
	invoiceItemSql "invoice-service/pkg/entities/invoice_items/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

type DBHandler struct {
	db                 *sharedDb.DbHandler
	logger             *logrus.Logger
	queries            *incomesql.Queries
	invoiceItemQueries *invoiceItemSql.Queries
}

func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := incomesql.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load income invoice SQL queries: %w", err)
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

func (h *DBHandler) Create(req *models.IncomeInvoiceCreateRequest) (*models.IncomeInvoice, error) {
	// Start transaction
	tx, err := h.db.BeginTx(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query, err := h.queries.Get(incomesql.CreateIncomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get create query: %w", err)
	}

	var invoice models.IncomeInvoice
	err = tx.QueryRow(query,
		req.OrderID,
		req.PaymentID,
		req.CustomerID,
		req.CustomerName,
		req.CustomerTaxID,
		req.InvoiceNumber,
		req.InvoiceType,
		req.Subtotal,
		req.TaxAmount,
		req.ServiceCharge,
		req.TotalAmount,
		req.PaymentMethod,
		req.XMLData,
		req.DigitalSignature,
		req.Status,
		req.GeneratedAt,
	).Scan(
		&invoice.ID, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to create income invoice")
		return nil, fmt.Errorf("failed to create income invoice: %w", err)
	}

	// Fill in the rest of the fields from the request
	invoice.OrderID = req.OrderID
	invoice.PaymentID = req.PaymentID
	invoice.CustomerID = req.CustomerID
	invoice.CustomerName = req.CustomerName
	invoice.CustomerTaxID = req.CustomerTaxID
	invoice.InvoiceNumber = req.InvoiceNumber
	invoice.InvoiceType = req.InvoiceType
	invoice.Subtotal = req.Subtotal
	invoice.TaxAmount = req.TaxAmount
	invoice.ServiceCharge = req.ServiceCharge
	invoice.TotalAmount = req.TotalAmount
	invoice.PaymentMethod = req.PaymentMethod
	invoice.XMLData = req.XMLData
	invoice.DigitalSignature = req.DigitalSignature
	invoice.Status = req.Status
	invoice.GeneratedAt = req.GeneratedAt

	// Create invoice items if provided
	if len(req.InvoiceItems) > 0 {
		for _, itemReq := range req.InvoiceItems {
			item, err := h.createInvoiceItem(tx, invoice.ID, "income", &itemReq)
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

	h.logger.WithField("id", invoice.ID).Info("Income invoice created with items")
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

func (h *DBHandler) GetByID(id string) (*models.IncomeInvoice, error) {
	query, err := h.queries.Get(incomesql.GetIncomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var invoice models.IncomeInvoice
	err = h.db.QueryRow(query, id).Scan(
		&invoice.ID,
		&invoice.OrderID,
		&invoice.PaymentID,
		&invoice.CustomerID,
		&invoice.CustomerName,
		&invoice.CustomerTaxID,
		&invoice.InvoiceNumber,
		&invoice.InvoiceType,
		&invoice.Subtotal,
		&invoice.TaxAmount,
		&invoice.ServiceCharge,
		&invoice.TotalAmount,
		&invoice.PaymentMethod,
		&invoice.XMLData,
		&invoice.DigitalSignature,
		&invoice.Status,
		&invoice.GeneratedAt,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("income invoice not found")
		}
		h.logger.WithError(err).Error("Failed to get income invoice")
		return nil, fmt.Errorf("failed to get income invoice: %w", err)
	}

	// Fetch invoice items for this invoice
	invoiceItems, err := h.getInvoiceItems(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}
	invoice.InvoiceItems = invoiceItems

	return &invoice, nil
}

func (h *DBHandler) Update(id string, req *models.IncomeInvoiceUpdateRequest) (*models.IncomeInvoice, error) {
	query, err := h.queries.Get(incomesql.UpdateIncomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get update query: %w", err)
	}

	_, err = h.db.Exec(query,
		req.PaymentID,
		req.CustomerID,
		req.CustomerName,
		req.CustomerTaxID,
		req.InvoiceType,
		req.Subtotal,
		req.TaxAmount,
		req.ServiceCharge,
		req.TotalAmount,
		req.PaymentMethod,
		req.XMLData,
		req.DigitalSignature,
		req.Status,
		req.GeneratedAt,
		id,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update income invoice")
		return nil, fmt.Errorf("failed to update income invoice: %w", err)
	}

	h.logger.WithField("id", id).Info("Income invoice updated")
	return h.GetByID(id)
}

func (h *DBHandler) Delete(id string) error {
	query, err := h.queries.Get(incomesql.DeleteIncomeInvoice)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(query, id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete income invoice")
		return fmt.Errorf("failed to delete income invoice: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("income invoice not found")
	}

	h.logger.WithField("id", id).Info("Income invoice deleted")
	return nil
}

func (h *DBHandler) List(req *models.IncomeInvoiceListRequest) (*models.IncomeInvoiceListResponse, error) {
	// Get the list and count queries
	listQuery, err := h.queries.Get(incomesql.ListIncomeInvoices)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	countQuery, err := h.queries.Get(incomesql.CountIncomeInvoices)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	// Get total count
	var total int
	err = h.db.QueryRow(countQuery, req.CustomerName, req.InvoiceType, req.Status, req.OrderID, req.CustomerID).Scan(&total)
	if err != nil {
		h.logger.WithError(err).Error("Failed to count income invoices")
		return nil, fmt.Errorf("failed to count income invoices: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	rows, err := h.db.Query(listQuery, req.CustomerName, req.InvoiceType, req.Status, req.OrderID, req.CustomerID, req.Limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list income invoices")
		return nil, fmt.Errorf("failed to list income invoices: %w", err)
	}
	defer rows.Close()

	var invoices []models.IncomeInvoice
	for rows.Next() {
		var invoice models.IncomeInvoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.OrderID,
			&invoice.PaymentID,
			&invoice.CustomerID,
			&invoice.CustomerName,
			&invoice.CustomerTaxID,
			&invoice.InvoiceNumber,
			&invoice.InvoiceType,
			&invoice.Subtotal,
			&invoice.TaxAmount,
			&invoice.ServiceCharge,
			&invoice.TotalAmount,
			&invoice.PaymentMethod,
			&invoice.XMLData,
			&invoice.DigitalSignature,
			&invoice.Status,
			&invoice.GeneratedAt,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			h.logger.WithError(err).Error("Failed to scan income invoice")
			return nil, fmt.Errorf("failed to scan income invoice: %w", err)
		}
		invoices = append(invoices, invoice)
	}

	return &models.IncomeInvoiceListResponse{
		Invoices: invoices,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}
