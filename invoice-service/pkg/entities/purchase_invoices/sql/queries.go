// Package sql provides SQL queries for purchase invoices
package sql

// Purchase invoice queries
const (
	CreatePurchaseInvoice = `
		INSERT INTO purchase_invoices (
			invoice_number, supplier_name, invoice_date, due_date,
			total_amount, status, image_url, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, invoice_number, supplier_name, invoice_date, due_date,
		          total_amount, status, image_url, notes, created_at, updated_at`

	GetPurchaseInvoiceByID = `
		SELECT id, invoice_number, supplier_name, invoice_date, due_date,
		       total_amount, status, image_url, notes, created_at, updated_at
		FROM purchase_invoices WHERE id = $1`

	UpdatePurchaseInvoice = `
		UPDATE purchase_invoices SET
			supplier_name = COALESCE($1, supplier_name),
			invoice_date = COALESCE($2, invoice_date),
			due_date = COALESCE($3, due_date),
			total_amount = COALESCE($4, total_amount),
			status = COALESCE($5, status),
			image_url = COALESCE($6, image_url),
			notes = COALESCE($7, notes),
			updated_at = NOW()
		WHERE id = $8`

	DeletePurchaseInvoice = `DELETE FROM purchase_invoices WHERE id = $1`

	ListPurchaseInvoices = `
		SELECT id, invoice_number, supplier_name, invoice_date, due_date,
		       total_amount, status, image_url, notes, created_at, updated_at
		FROM purchase_invoices
		WHERE ($1::text IS NULL OR supplier_name ILIKE '%' || $1 || '%')
		  AND ($2::text IS NULL OR status = $2)
		ORDER BY invoice_date DESC
		LIMIT $3 OFFSET $4`

	CountPurchaseInvoices = `
		SELECT COUNT(*) FROM purchase_invoices
		WHERE ($1::text IS NULL OR supplier_name ILIKE '%' || $1 || '%')
		  AND ($2::text IS NULL OR status = $2)`
)