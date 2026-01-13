package handlers

import (
	"database/sql"
	"fmt"

	"invoice-service/pkg/entities/income_invoices/models"
	incomesql "invoice-service/pkg/entities/income_invoices/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

type DBHandler struct {
	db      *sharedDb.DbHandler
	logger  *logrus.Logger
	queries *incomesql.Queries
}

func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := incomesql.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		logger:  logger,
		queries: queries,
	}, nil
}

func (h *DBHandler) Create(req *models.IncomeInvoiceCreateRequest) (*models.IncomeInvoice, error) {
	query, err := h.queries.Get(incomesql.CreateIncomeInvoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get create query: %w", err)
	}

	var invoice models.IncomeInvoice
	err = h.db.QueryRow(query,
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

	h.logger.WithField("id", invoice.ID).Info("Income invoice created")
	return &invoice, nil
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
