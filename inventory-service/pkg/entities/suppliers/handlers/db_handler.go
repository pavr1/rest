package handlers

import (
	"database/sql"
	"fmt"
	"inventory-service/pkg/entities/suppliers/models"
	supplierSQL "inventory-service/pkg/entities/suppliers/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for suppliers
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *supplierSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := supplierSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of suppliers
func (h *DBHandler) List(req *models.SupplierListRequest) (*models.SupplierListResponse, error) {
	offset := (req.Page - 1) * req.Limit

	countQuery, err := h.queries.Get(supplierSQL.CountSuppliersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, req.Name, req.Email, req.Phone).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count suppliers: %w", err)
	}

	listQuery, err := h.queries.Get(supplierSQL.ListSuppliersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, req.Name, req.Email, req.Phone, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var supplier models.Supplier
		err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.ContactName,
			&supplier.Phone,
			&supplier.Email,
			&supplier.Address,
			&supplier.CreatedAt,
			&supplier.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan supplier: %w", err)
		}
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating suppliers: %w", err)
	}

	return &models.SupplierListResponse{
		Suppliers: suppliers,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
	}, nil
}

// GetByID retrieves a supplier by ID
func (h *DBHandler) GetByID(id string) (*models.Supplier, error) {
	query, err := h.queries.Get(supplierSQL.GetSupplierByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var supplier models.Supplier
	err = h.db.QueryRow(query, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.ContactName,
		&supplier.Phone,
		&supplier.Email,
		&supplier.Address,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	return &supplier, nil
}

// Create creates a new supplier
func (h *DBHandler) Create(req *models.SupplierCreateRequest) (*models.Supplier, error) {
	query, err := h.queries.Get(supplierSQL.CreateSupplierQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var supplier models.Supplier
	err = h.db.QueryRow(query, req.Name, req.ContactName, req.Phone, req.Email, req.Address).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.ContactName,
		&supplier.Phone,
		&supplier.Email,
		&supplier.Address,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"supplier_id": supplier.ID,
		"name":        supplier.Name,
	}).Info("Supplier created successfully")

	return &supplier, nil
}

// Update updates an existing supplier
func (h *DBHandler) Update(id string, req *models.SupplierUpdateRequest) (*models.Supplier, error) {
	query, err := h.queries.Get(supplierSQL.UpdateSupplierQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var supplier models.Supplier
	err = h.db.QueryRow(query, id, req.Name, req.ContactName, req.Phone, req.Email, req.Address).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.ContactName,
		&supplier.Phone,
		&supplier.Email,
		&supplier.Address,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"supplier_id": supplier.ID,
		"name":        supplier.Name,
	}).Info("Supplier updated successfully")

	return &supplier, nil
}

// Delete deletes a supplier
func (h *DBHandler) Delete(id string) error {
	// Check for dependencies
	deps, err := h.checkDependencies(id)
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if deps.PurchaseInvoiceCount > 0 || deps.OutcomeInvoiceCount > 0 {
		return fmt.Errorf("cannot delete supplier: it has %d purchase invoices and %d outcome invoices",
			deps.PurchaseInvoiceCount, deps.OutcomeInvoiceCount)
	}

	query, err := h.queries.Get(supplierSQL.DeleteSupplierQuery)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	result, err := h.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier not found")
	}

	h.logger.WithField("supplier_id", id).Info("Supplier deleted successfully")
	return nil
}

// checkDependencies checks if a supplier has dependencies
func (h *DBHandler) checkDependencies(id string) (*SupplierDependencies, error) {
	query, err := h.queries.Get(supplierSQL.CheckSupplierDependenciesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var deps SupplierDependencies
	err = h.db.QueryRow(query, id).Scan(&deps.PurchaseInvoiceCount, &deps.OutcomeInvoiceCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check dependencies: %w", err)
	}

	return &deps, nil
}

// SupplierDependencies represents the dependencies of a supplier
type SupplierDependencies struct {
	PurchaseInvoiceCount int `json:"purchase_invoice_count"`
	OutcomeInvoiceCount  int `json:"outcome_invoice_count"`
}