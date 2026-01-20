-- List invoice items for a specific invoice
SELECT
    ii.id,
    ii.invoice_id,
    ii.stock_variant_id,
    ii.detail,
    ii.count,
    ii.unit_type,
    ii.price,
    ii.items_per_unit,
    ii.total,
    ii.expiration_date,
    ii.created_at,
    ii.updated_at,
    sv.name AS stock_variant_name
FROM invoice_items ii
LEFT JOIN stock_variants sv ON ii.stock_variant_id = sv.id
WHERE ii.invoice_id = $1
ORDER BY ii.created_at DESC;
