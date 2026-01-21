-- Calculate and update the average cost per portion for a stock variant
-- Only considers active stock counts (is_out = false) with cost_per_portion > 0
UPDATE stock_variants
SET avg_cost = COALESCE(
    (SELECT AVG(cost_per_portion) 
     FROM stock_count 
     WHERE stock_variant_id = $1 
       AND is_out = false 
       AND cost_per_portion IS NOT NULL 
       AND cost_per_portion > 0),
    0
),
updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, avg_cost;
