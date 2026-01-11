-- Migration: Add Sub Menus table and update Menu Items
-- This creates a 3-tier menu hierarchy: Menu Categories → Sub Menus → Menu Items

-- =============================================================================
-- Step 1: Create Sub Menus table
-- =============================================================================
CREATE TABLE sub_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID NOT NULL REFERENCES menu_categories(id) ON DELETE RESTRICT,
    image_url VARCHAR(500),
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('kitchen', 'bar')),
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- Step 2: Create indexes for sub_menus
-- =============================================================================
CREATE INDEX idx_sub_menus_category ON sub_menus(category_id);
CREATE INDEX idx_sub_menus_active ON sub_menus(is_active);
CREATE INDEX idx_sub_menus_item_type ON sub_menus(item_type);

-- =============================================================================
-- Step 3: Create trigger for sub_menus updated_at
-- =============================================================================
CREATE TRIGGER update_sub_menus_updated_at BEFORE UPDATE ON sub_menus 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- Step 4: Add sub_menu_id column to menu_items (nullable for now)
-- =============================================================================
ALTER TABLE menu_items ADD COLUMN sub_menu_id UUID;

-- =============================================================================
-- Step 5: Migrate existing menu_items to sub_menus (if any exist)
-- Strategy: Create one sub_menu per existing menu_item (using menu_item name)
-- This preserves the existing structure - each old menu_item becomes a sub_menu
-- Then the old menu_items become the new detailed menu_items
-- =============================================================================
INSERT INTO sub_menus (id, name, description, category_id, image_url, item_type, display_order, is_active, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    mi.name,
    mi.description,
    mi.category_id,
    mi.image_url,
    mi.item_type,
    0,
    mi.is_available,
    mi.created_at,
    mi.updated_at
FROM menu_items mi;

-- =============================================================================
-- Step 6: Update menu_items to reference the migrated sub_menus
-- Match by name, category_id, and item_type (1:1 mapping)
-- =============================================================================
UPDATE menu_items mi
SET sub_menu_id = sm.id
FROM sub_menus sm
WHERE mi.name = sm.name 
  AND mi.category_id = sm.category_id 
  AND mi.item_type = sm.item_type;

-- =============================================================================
-- Step 7: Add preparation_time and display_order columns to menu_items
-- =============================================================================
ALTER TABLE menu_items ADD COLUMN preparation_time INTEGER;
ALTER TABLE menu_items ADD COLUMN display_order INTEGER NOT NULL DEFAULT 0;

-- =============================================================================
-- Step 8: Drop old columns from menu_items that are now in sub_menus
-- =============================================================================
-- Drop the old category_id column (now in sub_menus)
ALTER TABLE menu_items DROP COLUMN category_id;
-- Drop item_type (now inherited from sub_menu)
ALTER TABLE menu_items DROP COLUMN item_type;

-- =============================================================================
-- Step 9: Make sub_menu_id NOT NULL after data migration
-- PostgreSQL allows setting NOT NULL on empty tables (just enforces for future inserts)
-- =============================================================================
ALTER TABLE menu_items ALTER COLUMN sub_menu_id SET NOT NULL;

-- =============================================================================
-- Step 10: Drop old indexes and create new ones for menu_items
-- =============================================================================
DROP INDEX IF EXISTS idx_menu_items_category;
DROP INDEX IF EXISTS idx_menu_items_type;
CREATE INDEX idx_menu_items_sub_menu ON menu_items(sub_menu_id);

-- =============================================================================
-- Step 11: Update menu_item_stock_items to use correct constraint names if needed
-- (The FK to menu_items should still work since menu_items.id didn't change)
-- =============================================================================
-- No changes needed - menu_items.id is still the PK

