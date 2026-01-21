-- Migration 007: Rollback Stock Count Schema Changes

-- Restore dropped columns on stock_variants
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS invoice_id UUID;
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS unit VARCHAR(50);
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS number_of_units DECIMAL(10,2);

-- Set defaults for restored columns
UPDATE stock_variants SET unit = 'Unit' WHERE unit IS NULL;
UPDATE stock_variants SET number_of_units = 1 WHERE number_of_units IS NULL;

-- Add constraints back
ALTER TABLE stock_variants ALTER COLUMN unit SET NOT NULL;
ALTER TABLE stock_variants ALTER COLUMN number_of_units SET NOT NULL;
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'stock_variants_number_of_units_check') THEN
        ALTER TABLE stock_variants ADD CONSTRAINT stock_variants_number_of_units_check CHECK (number_of_units > 0);
    END IF;
END $$;

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
