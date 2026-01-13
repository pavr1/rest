-- Update an income invoice
UPDATE income_invoices SET
    payment_id = COALESCE($2, payment_id),
    customer_id = COALESCE($3, customer_id),
    customer_name = COALESCE($4, customer_name),
    customer_tax_id = COALESCE($5, customer_tax_id),
    invoice_type = COALESCE($6, invoice_type),
    subtotal = COALESCE($7, subtotal),
    tax_amount = COALESCE($8, tax_amount),
    service_charge = COALESCE($9, service_charge),
    total_amount = COALESCE($10, total_amount),
    payment_method = COALESCE($11, payment_method),
    xml_data = COALESCE($12, xml_data),
    digital_signature = COALESCE($13, digital_signature),
    status = COALESCE($14, status),
    generated_at = COALESCE($15, generated_at),
    updated_at = NOW()
WHERE id = $1
RETURNING id, updated_at;