-- Migration 007: Stock Count Schema Changes
-- Purpose: Add stock_count table, simplify stock_variants, update invoice_items and menu_ingredients
-- Note: This migration handles both fresh installs (where base schema is already updated) 
-- and upgrades from older schema versions.

-- Step 1: Create stock_count table (IF NOT EXISTS for fresh install compatibility)
CREATE TABLE IF NOT EXISTS stock_count (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_variant_id UUID NOT NULL REFERENCES stock_variants(id) ON DELETE CASCADE,
    invoice_id UUID NOT NULL REFERENCES outcome_invoices(id) ON DELETE CASCADE,
    count DECIMAL(10,2) NOT NULL CHECK (count > 0),
    unit VARCHAR(50) NOT NULL,
    purchased_at TIMESTAMP NOT NULL,
    is_out BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 2: Create indexes for stock_count (IF NOT EXISTS)
CREATE INDEX IF NOT EXISTS idx_stock_count_variant ON stock_count(stock_variant_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_invoice ON stock_count(invoice_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_is_out ON stock_count(stock_variant_id, is_out);

-- Step 3: Create trigger for stock_count updated_at (only if not exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_stock_count_updated_at') THEN
        CREATE TRIGGER update_stock_count_updated_at BEFORE UPDATE ON stock_count
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Step 4: Add description column to stock_variants (if not exists)
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS description TEXT;

-- Step 4b: Remove deprecated columns from stock_variants (if they exist)
ALTER TABLE stock_variants DROP CONSTRAINT IF EXISTS stock_variants_number_of_units_check;
ALTER TABLE stock_variants DROP COLUMN IF EXISTS unit;
ALTER TABLE stock_variants DROP COLUMN IF EXISTS number_of_units;
ALTER TABLE stock_variants DROP COLUMN IF EXISTS invoice_id;

-- Step 5: Add stock_variant_id to invoice_items (if not exists)
ALTER TABLE invoice_items ADD COLUMN IF NOT EXISTS stock_variant_id UUID REFERENCES stock_variants(id) ON DELETE SET NULL;

-- Step 6: Create index for invoice_items.stock_variant_id (IF NOT EXISTS)
CREATE INDEX IF NOT EXISTS idx_invoice_items_stock_variant ON invoice_items(stock_variant_id);

-- Step 7: Update menu_ingredients to use stock_variant_id (if not exists)
ALTER TABLE menu_ingredients ADD COLUMN IF NOT EXISTS stock_variant_id UUID REFERENCES stock_variants(id) ON DELETE RESTRICT;

-- Step 8: Create index for menu_ingredients.stock_variant_id (IF NOT EXISTS)
CREATE INDEX IF NOT EXISTS idx_menu_ingredients_stock_variant ON menu_ingredients(stock_variant_id);
