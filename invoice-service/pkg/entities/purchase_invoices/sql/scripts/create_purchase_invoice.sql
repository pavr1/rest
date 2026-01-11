INSERT INTO purchase_invoices (
    invoice_number, supplier_name, invoice_date, due_date,
    total_amount, status, image_url, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, invoice_number, supplier_name, invoice_date, due_date,
          total_amount, status, image_url, notes, created_at, updated_at
