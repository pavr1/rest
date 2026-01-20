SELECT 
    sc.id, 
    sc.stock_variant_id, 
    sc.invoice_id, 
    sc.count, 
    sc.unit, 
    sc.purchased_at, 
    sc.is_out, 
    sc.created_at, 
    sc.updated_at,
    sv.name AS stock_variant_name,
    oi.invoice_number,
    s.name AS supplier_name
FROM stock_count sc
LEFT JOIN stock_variants sv ON sc.stock_variant_id = sv.id
LEFT JOIN outcome_invoices oi ON sc.invoice_id = oi.id
LEFT JOIN suppliers s ON oi.supplier_id = s.id
ORDER BY sc.purchased_at DESC
LIMIT $1 OFFSET $2;
