// Package sql provides SQL queries for existences
package sql

// Existence queries
const (
	CreateExistence = `
		INSERT INTO existences (
			invoice_detail_id, stock_item_id, units_purchased, cost_per_unit,
			total_cost, expiry_date, batch_number, current_stock
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, invoice_detail_id, stock_item_id, units_purchased,
		          cost_per_unit, total_cost, expiry_date, batch_number,
		          current_stock, created_at, updated_at`

	GetExistenceByID = `
		SELECT e.id, e.invoice_detail_id, e.stock_item_id, si.name as stock_item_name,
		       e.units_purchased, e.cost_per_unit, e.total_cost, e.expiry_date,
		       e.batch_number, e.current_stock, e.created_at, e.updated_at
		FROM existences e
		LEFT JOIN stock_items si ON e.stock_item_id = si.id
		WHERE e.id = $1`

	UpdateExistence = `
		UPDATE existences SET
			current_stock = COALESCE($1, current_stock),
			expiry_date = COALESCE($2, expiry_date),
			batch_number = COALESCE($3, batch_number),
			updated_at = NOW()
		WHERE id = $4`

	DeleteExistence = `DELETE FROM existences WHERE id = $1`

	ListExistences = `
		SELECT e.id, e.invoice_detail_id, e.stock_item_id, si.name as stock_item_name,
		       e.units_purchased, e.cost_per_unit, e.total_cost, e.expiry_date,
		       e.batch_number, e.current_stock, e.created_at, e.updated_at
		FROM existences e
		LEFT JOIN stock_items si ON e.stock_item_id = si.id
		WHERE ($1::text IS NULL OR e.stock_item_id = $1)
		ORDER BY e.created_at DESC
		LIMIT $2 OFFSET $3`

	CountExistences = `
		SELECT COUNT(*) FROM existences e
		WHERE ($1::text IS NULL OR e.stock_item_id = $1)`

	GetInvoiceDetailForExistence = `
		SELECT id.stock_item_id, si.name, id.quantity, id.unit_price,
		       id.unit_of_measure, id.items_per_unit, id.expiry_date, id.batch_number
		FROM invoice_details id
		LEFT JOIN stock_items si ON id.stock_item_id = si.id
		WHERE id.id = $1`

	UpdateStockItemFromExistence = `
		UPDATE stock_items
		SET current_stock = $1, unit_cost = $2, total_value = $3, updated_at = NOW()
		WHERE id = $4`
)