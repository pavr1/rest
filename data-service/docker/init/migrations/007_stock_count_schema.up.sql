-- Migration 007: Stock Count Schema Changes
-- Purpose: Add stock_count table, simplify stock_variants, update invoice_items and menu_ingredients

-- Step 1: Create stock_count table
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

-- Step 2: Create indexes for stock_count
CREATE INDEX IF NOT EXISTS idx_stock_count_variant ON stock_count(stock_variant_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_invoice ON stock_count(invoice_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_is_out ON stock_count(stock_variant_id, is_out);

-- Step 3: Create trigger for stock_count updated_at
CREATE TRIGGER update_stock_count_updated_at BEFORE UPDATE ON stock_count
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Step 4: Add description column to stock_variants (if not exists)
ALTER TABLE stock_variants ADD COLUMN IF NOT EXISTS description TEXT;

-- Step 4b: Make deprecated stock_variants columns nullable
ALTER TABLE stock_variants ALTER COLUMN unit DROP NOT NULL;
ALTER TABLE stock_variants ALTER COLUMN number_of_units DROP NOT NULL;
ALTER TABLE stock_variants DROP CONSTRAINT IF EXISTS stock_variants_number_of_units_check;

-- Step 5: Add stock_variant_id to invoice_items
ALTER TABLE invoice_items ADD COLUMN IF NOT EXISTS stock_variant_id UUID REFERENCES stock_variants(id) ON DELETE SET NULL;

-- Step 6: Create index for invoice_items.stock_variant_id
CREATE INDEX IF NOT EXISTS idx_invoice_items_stock_variant ON invoice_items(stock_variant_id);

-- Step 7: Update menu_ingredients to use stock_variant_id instead of stock_sub_category_id
-- First add the new column
ALTER TABLE menu_ingredients ADD COLUMN IF NOT EXISTS stock_variant_id UUID REFERENCES stock_variants(id) ON DELETE RESTRICT;

-- Step 8: Create index for menu_ingredients.stock_variant_id
CREATE INDEX IF NOT EXISTS idx_menu_ingredients_stock_variant ON menu_ingredients(stock_variant_id);

-- Note: The following columns are being deprecated but kept for backward compatibility:
-- - stock_variants.unit (can be dropped later after data migration)
-- - stock_variants.number_of_units (can be dropped later after data migration)
-- - stock_variants.invoice_id (can be dropped later after data migration)
-- - invoice_items.inventory_category_id (can be dropped later after data migration)
-- - invoice_items.inventory_sub_category_id (can be dropped later after data migration)
-- - menu_ingredients.stock_sub_category_id (can be dropped later after data migration)

-- Step 9: Comment on deprecated columns
COMMENT ON COLUMN stock_variants.unit IS 'DEPRECATED: Use stock_count.unit instead';
COMMENT ON COLUMN stock_variants.number_of_units IS 'DEPRECATED: Use stock_count.count instead';
COMMENT ON COLUMN stock_variants.invoice_id IS 'DEPRECATED: Use stock_count.invoice_id instead';
COMMENT ON COLUMN invoice_items.inventory_category_id IS 'DEPRECATED: Use stock_variant_id instead';
COMMENT ON COLUMN invoice_items.inventory_sub_category_id IS 'DEPRECATED: Use stock_variant_id instead';
COMMENT ON COLUMN menu_ingredients.stock_sub_category_id IS 'DEPRECATED: Use stock_variant_id instead';
