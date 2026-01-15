-- Update an income invoice
UPDATE income_invoices SET
    payment_id = COALESCE($2, payment_id),
    customer_id = COALESCE($3, customer_id),
    invoice_type = COALESCE($4, invoice_type),
    subtotal = COALESCE($5, subtotal),
    tax_amount = COALESCE($6, tax_amount),
    service_charge = COALESCE($7, service_charge),
    total_amount = COALESCE($8, total_amount),
    payment_method = COALESCE($9, payment_method),
    xml_data = COALESCE($10, xml_data),
    digital_signature = COALESCE($11, digital_signature),
    status = COALESCE($12, status),
    generated_at = COALESCE($13, generated_at),
    updated_at = NOW()
WHERE id = $1
RETURNING id, updated_at;
