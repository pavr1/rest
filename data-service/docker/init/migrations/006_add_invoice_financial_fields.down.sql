-- Remove financial/stock variant fields and restore old structure
-- Remove financial fields from outcome_invoices
ALTER TABLE outcome_invoices
DROP COLUMN IF EXISTS discount_amount,
DROP COLUMN IF EXISTS tax_amount,
DROP COLUMN IF EXISTS subtotal,
DROP COLUMN IF EXISTS due_date;

-- Remove invoice_id from stock_variants
ALTER TABLE stock_variants
DROP COLUMN IF EXISTS invoice_id;

-- Restore removed columns to invoice_items
ALTER TABLE invoice_items
ADD COLUMN IF NOT EXISTS stock_variant_id UUID;

-- Restore inventory fields to outcome_invoices
ALTER TABLE outcome_invoices
ADD COLUMN IF NOT EXISTS inventory_category_id UUID REFERENCES stock_categories(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS inventory_sub_category_id UUID REFERENCES stock_sub_categories(id) ON DELETE SET NULL;
