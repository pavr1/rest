// Package sql provides SQL queries for invoice details
package sql

// Invoice detail queries
const (
	CreateInvoiceDetail = `
		INSERT INTO invoice_details (
			invoice_id, stock_item_id, description, quantity, unit_of_measure,
			items_per_unit, unit_price, total_price, expiry_date, batch_number
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, invoice_id, stock_item_id, description, quantity,
		          unit_of_measure, items_per_unit, unit_price, total_price,
		          expiry_date, batch_number, created_at, updated_at`

	GetInvoiceDetailByID = `
		SELECT id, invoice_id, stock_item_id, description, quantity, unit_of_measure,
		       items_per_unit, unit_price, total_price, expiry_date, batch_number,
		       created_at, updated_at
		FROM invoice_details WHERE id = $1 AND invoice_id = $2`

	UpdateInvoiceDetail = `
		UPDATE invoice_details SET
			stock_item_id = COALESCE($1, stock_item_id),
			description = COALESCE($2, description),
			quantity = COALESCE($3, quantity),
			unit_of_measure = COALESCE($4, unit_of_measure),
			items_per_unit = COALESCE($5, items_per_unit),
			unit_price = COALESCE($6, unit_price),
			total_price = COALESCE($7, total_price),
			expiry_date = COALESCE($8, expiry_date),
			batch_number = COALESCE($9, batch_number),
			updated_at = NOW()
		WHERE id = $10 AND invoice_id = $11`

	DeleteInvoiceDetail = `DELETE FROM invoice_details WHERE id = $1 AND invoice_id = $2`

	ListInvoiceDetails = `
		SELECT id, invoice_id, stock_item_id, si.name as stock_item_name,
		       description, quantity, unit_of_measure, items_per_unit,
		       unit_price, total_price, expiry_date, batch_number,
		       id.created_at, id.updated_at
		FROM invoice_details id
		LEFT JOIN stock_items si ON id.stock_item_id = si.id
		WHERE id.invoice_id = $1
		ORDER BY id.created_at`
)