-- Rollback: Remove Sub Menus table and revert Menu Items structure

-- =============================================================================
-- Step 1: Add back columns to menu_items
-- =============================================================================
ALTER TABLE menu_items ADD COLUMN category_id UUID;
ALTER TABLE menu_items ADD COLUMN item_type VARCHAR(20);

-- =============================================================================
-- Step 2: Restore data from sub_menus back to menu_items
-- =============================================================================
UPDATE menu_items mi
SET 
    category_id = sm.category_id,
    item_type = sm.item_type
FROM sub_menus sm
WHERE mi.sub_menu_id = sm.id;

-- =============================================================================
-- Step 3: Make restored columns NOT NULL and add constraints
-- =============================================================================
ALTER TABLE menu_items ALTER COLUMN category_id SET NOT NULL;
ALTER TABLE menu_items ADD CONSTRAINT menu_items_category_id_fkey 
    FOREIGN KEY (category_id) REFERENCES menu_categories(id) ON DELETE RESTRICT;
ALTER TABLE menu_items ALTER COLUMN item_type SET NOT NULL;
ALTER TABLE menu_items ADD CONSTRAINT menu_items_item_type_check 
    CHECK (item_type IN ('kitchen', 'bar'));

-- =============================================================================
-- Step 4: Drop new columns from menu_items
-- =============================================================================
ALTER TABLE menu_items DROP COLUMN sub_menu_id;
ALTER TABLE menu_items DROP COLUMN preparation_time;
ALTER TABLE menu_items DROP COLUMN display_order;

-- =============================================================================
-- Step 5: Recreate old indexes
-- =============================================================================
DROP INDEX IF EXISTS idx_menu_items_sub_menu;
CREATE INDEX idx_menu_items_category ON menu_items(category_id);
CREATE INDEX idx_menu_items_type ON menu_items(item_type);

-- =============================================================================
-- Step 6: Drop sub_menus trigger
-- =============================================================================
DROP TRIGGER IF EXISTS update_sub_menus_updated_at ON sub_menus;

-- =============================================================================
-- Step 7: Drop sub_menus indexes
-- =============================================================================
DROP INDEX IF EXISTS idx_sub_menus_category;
DROP INDEX IF EXISTS idx_sub_menus_active;
DROP INDEX IF EXISTS idx_sub_menus_item_type;

-- =============================================================================
-- Step 8: Drop sub_menus table
-- =============================================================================
DROP TABLE IF EXISTS sub_menus;
