-- List income invoices with optional filters
SELECT
    id,
    order_id,
    payment_id,
    customer_id,
    invoice_number,
    invoice_type,
    subtotal,
    tax_amount,
    service_charge,
    total_amount,
    payment_method,
    xml_data,
    digital_signature,
    status,
    generated_at,
    created_at,
    updated_at
FROM income_invoices
WHERE ($1::text IS NULL OR customer_id ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR invoice_type = $2)
  AND ($3::text IS NULL OR status = $3)
  AND ($4::text IS NULL OR order_id = $4)
ORDER BY created_at DESC
LIMIT $5 OFFSET $6;
