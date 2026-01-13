package handlers

import (
	"database/sql"
	"fmt"

	"invoice-service/pkg/entities/invoice_items/models"
	invoiceitemsql "invoice-service/pkg/entities/invoice_items/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

type DBHandler struct {
	db      *sharedDb.DbHandler
	logger  *logrus.Logger
	queries *invoiceitemsql.Queries
}

func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := invoiceitemsql.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		logger:  logger,
		queries: queries,
	}, nil
}

func (h *DBHandler) Create(invoiceID string, invoiceType string, req *models.InvoiceItemCreateRequest) (*models.InvoiceItem, error) {
	// Calculate total
	total := req.Count * req.Price

	query, err := h.queries.Get(invoiceitemsql.CreateInvoiceItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get create query: %w", err)
	}

	var item models.InvoiceItem
	err = h.db.QueryRow(query,
		invoiceID, invoiceType, req.Detail, req.Count, req.UnitType,
		req.Price, req.ItemsPerUnit, req.ExpirationDate,
	).Scan(
		&item.ID, &item.Total, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to create invoice item")
		return nil, fmt.Errorf("failed to create invoice item: %w", err)
	}

	// Fill in the rest of the fields from the request
	item.InvoiceID = invoiceID
	item.Detail = req.Detail
	item.Count = req.Count
	item.UnitType = req.UnitType
	item.Price = req.Price
	item.ItemsPerUnit = req.ItemsPerUnit
	item.Total = total
	item.ExpirationDate = req.ExpirationDate

	h.logger.WithField("id", item.ID).Info("Invoice item created")
	return &item, nil
}

func (h *DBHandler) GetByID(id string) (*models.InvoiceItem, error) {
	query, err := h.queries.Get(invoiceitemsql.GetInvoiceItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var item models.InvoiceItem
	err = h.db.QueryRow(query, id).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice item not found")
		}
		h.logger.WithError(err).Error("Failed to get invoice item")
		return nil, fmt.Errorf("failed to get invoice item: %w", err)
	}

	return &item, nil
}

func (h *DBHandler) Update(id string, req *models.InvoiceItemUpdateRequest) (*models.InvoiceItem, error) {
	// Calculate new total if count or price changed
	var total *float64
	if req.Count != nil || req.Price != nil {
		current, err := h.GetByID(id)
		if err != nil {
			return nil, err
		}

		newCount := current.Count
		newPrice := current.Price

		if req.Count != nil {
			newCount = *req.Count
		}
		if req.Price != nil {
			newPrice = *req.Price
		}

		calculatedTotal := newCount * newPrice
		total = &calculatedTotal
	}

	query, err := h.queries.Get(invoiceitemsql.UpdateInvoiceItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get update query: %w", err)
	}

	_, err = h.db.Exec(query,
		req.Detail,
		req.Count,
		req.UnitType,
		req.Price,
		req.ItemsPerUnit,
		total,
		req.ExpirationDate,
		id,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update invoice item")
		return nil, fmt.Errorf("failed to update invoice item: %w", err)
	}

	h.logger.WithField("id", id).Info("Invoice item updated")
	return h.GetByID(id)
}

func (h *DBHandler) Delete(id string) error {
	query, err := h.queries.Get(invoiceitemsql.DeleteInvoiceItem)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(query, id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete invoice item")
		return fmt.Errorf("failed to delete invoice item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invoice item not found")
	}

	h.logger.WithField("id", id).Info("Invoice item deleted")
	return nil
}

func (h *DBHandler) ListByInvoice(invoiceID string) (*models.InvoiceItemListResponse, error) {
	query, err := h.queries.Get(invoiceitemsql.ListInvoiceItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(query, invoiceID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list invoice items")
		return nil, fmt.Errorf("failed to list invoice items: %w", err)
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

	return &models.InvoiceItemListResponse{
		Items: items,
		Total: len(items),
	}, nil
}
