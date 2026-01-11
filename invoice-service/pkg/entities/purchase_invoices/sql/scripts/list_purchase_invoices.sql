SELECT id, invoice_number, supplier_name, invoice_date, due_date,
       total_amount, status, image_url, notes, created_at, updated_at
FROM purchase_invoices
WHERE ($1::text IS NULL OR supplier_name ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR status = $2)
ORDER BY invoice_date DESC
LIMIT $3 OFFSET $4
