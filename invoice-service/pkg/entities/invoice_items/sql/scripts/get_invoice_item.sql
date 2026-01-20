-- Get invoice item by ID
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
WHERE ii.id = $1;
