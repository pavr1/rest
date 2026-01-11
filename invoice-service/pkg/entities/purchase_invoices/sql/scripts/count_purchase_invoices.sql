SELECT COUNT(*) FROM purchase_invoices
WHERE ($1::text IS NULL OR supplier_name ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR status = $2)
