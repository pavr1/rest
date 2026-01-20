-- Migration 007: Rollback Stock Count Schema Changes

-- Remove comments on deprecated columns
COMMENT ON COLUMN invoice_items.inventory_category_id IS NULL;
COMMENT ON COLUMN invoice_items.inventory_sub_category_id IS NULL;
COMMENT ON COLUMN menu_ingredients.stock_sub_category_id IS NULL;

-- Restore dropped columns on stock_variants
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS invoice_id UUID;
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS unit VARCHAR(50) NOT NULL DEFAULT 'Unit';
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS number_of_units DECIMAL(10,2) NOT NULL DEFAULT 1;
ALTER TABLE stock_variants ADD CONSTRAINT stock_variants_number_of_units_check CHECK (number_of_units > 0);

-- Drop new columns
ALTER TABLE menu_ingredients DROP COLUMN IF EXISTS stock_variant_id;
ALTER TABLE invoice_items DROP COLUMN IF EXISTS stock_variant_id;
ALTER TABLE stock_variants DROP COLUMN IF EXISTS description;

-- Drop indexes
DROP INDEX IF EXISTS idx_menu_ingredients_stock_variant;
DROP INDEX IF EXISTS idx_invoice_items_stock_variant;
DROP INDEX IF EXISTS idx_stock_count_is_out;
DROP INDEX IF EXISTS idx_stock_count_invoice;
DROP INDEX IF EXISTS idx_stock_count_variant;

-- Drop trigger
DROP TRIGGER IF EXISTS update_stock_count_updated_at ON stock_count;

-- Drop stock_count table
DROP TABLE IF EXISTS stock_count;
