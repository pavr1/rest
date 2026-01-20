package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"invoice-service/pkg/entities/income_invoices/models"
	incomesql "invoice-service/pkg/entities/income_invoices/sql"
	invoiceItemModels "invoice-service/pkg/entities/invoice_items/models"
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

// Create creates a new income invoice with its items in a transaction
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
	//pvillalobos TODO: check if the invoice items are empty
	if len(req.InvoiceItems) > 0 {
		for _, itemReq := range req.InvoiceItems {
			itemReq.InvoiceID = invoice.ID

			item, err := h.createInvoiceItem(tx, &itemReq)
			if err != nil {
				return nil, fmt.Errorf("failed to create invoice item: %w", err)
			}
			invoice.InvoiceItems = append(invoice.InvoiceItems, *item)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &invoice, nil
}

// GetByID retrieves an income invoice by ID with its items
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

	// Get invoice items
	invoiceItems, err := h.getInvoiceItems(invoice.ID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get invoice items")
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}
	invoice.InvoiceItems = invoiceItems

	return &invoice, nil
}

// Update updates an income invoice (transaction support can be added if items need updating)
func (h *DBHandler) Update(id string, req *models.IncomeInvoiceUpdateRequest) (*models.IncomeInvoice, error) {
	query, err := h.queries.Get(incomesql.UpdateIncomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get update query: %w", err)
	}

	_, err = h.db.Exec(query,
		id,
		req.PaymentID,
		req.CustomerID,
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
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to update income invoice")
		return nil, fmt.Errorf("failed to update income invoice: %w", err)
	}

	// Return updated invoice
	return h.GetByID(id)
}

// Delete deletes an income invoice (transaction support can be added if cascading deletes are needed)
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

	return nil
}

// List retrieves income invoices with pagination and filtering
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
	err = h.db.QueryRow(countQuery, req.CustomerID, req.InvoiceType, req.Status, req.OrderID).Scan(&total)
	if err != nil {
		h.logger.WithError(err).Error("Failed to count income invoices")
		return nil, fmt.Errorf("failed to count income invoices: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	rows, err := h.db.Query(listQuery, req.CustomerID, req.InvoiceType, req.Status, req.OrderID, req.Limit, offset)
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

		// Get invoice items for this invoice
		invoiceItems, err := h.getInvoiceItems(invoice.ID)
		if err != nil {
			h.logger.WithError(err).Error("Failed to get invoice items")
			return nil, fmt.Errorf("failed to get invoice items: %w", err)
		}
		invoice.InvoiceItems = invoiceItems

		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &models.IncomeInvoiceListResponse{
		Invoices: invoices,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

// createInvoiceItem creates a single invoice item within a transaction
func (h *DBHandler) createInvoiceItem(tx *sql.Tx, req *invoiceItemModels.InvoiceItemCreateRequest) (*invoiceItemModels.InvoiceItem, error) {
	query, err := h.invoiceItemQueries.Get(invoiceItemSql.CreateInvoiceItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get create item query: %w", err)
	}

	var item invoiceItemModels.InvoiceItem
	err = tx.QueryRow(query,
		req.InvoiceID,
		req.Detail,
		req.Count,
		req.UnitType,
		req.Price,
		req.ItemsPerUnit,
		req.ExpirationDate,
	).Scan(
		&item.ID, &item.Total, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create invoice item: %w", err)
	}

	// Fill in the rest of the fields
	item.InvoiceID = req.InvoiceID
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
		return nil, fmt.Errorf("failed to get list items query: %w", err)
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
		var inventoryCategoryID sql.NullString
		var inventorySubCategoryID sql.NullString
		var detail string
		err := rows.Scan(
			&item.ID,
			&item.InvoiceID,
			&inventoryCategoryID,
			&inventorySubCategoryID,
			&detail,
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
			return nil, fmt.Errorf("failed to scan invoice item: %w", err)
		}

		if inventoryCategoryID.Valid {
			item.InventoryCategoryID = &inventoryCategoryID.String
		}

		if inventorySubCategoryID.Valid {
			item.InventorySubCategoryID = &inventorySubCategoryID.String
		}

		item.Detail = detail

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating item rows: %w", err)
	}

	return items, nil
}
