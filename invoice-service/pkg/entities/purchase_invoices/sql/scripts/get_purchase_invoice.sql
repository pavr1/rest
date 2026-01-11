SELECT id, invoice_number, supplier_name, invoice_date, due_date,
       total_amount, status, image_url, notes, created_at, updated_at
FROM purchase_invoices WHERE id = $1
