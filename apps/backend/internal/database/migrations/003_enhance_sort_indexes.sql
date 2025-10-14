-- +migrate Up
-- --- Enhancements for sorting, indexing, and data integrity ---

-- Add explicit sort_order columns
ALTER TABLE plants ADD COLUMN sort_order SERIAL;
ALTER TABLE observations ADD COLUMN sort_order SERIAL;

-- Add composite indexes for efficient sorting
CREATE INDEX idx_plants_user_sort ON plants(user_id, sort_order);
CREATE INDEX idx_observations_plant_sort ON observations(plant_id, sort_order);
CREATE INDEX idx_observations_user_plant_date_desc
    ON observations(user_id, plant_id, date DESC);

-- Prevent invalid plant measurements
ALTER TABLE observations ADD CONSTRAINT height_nonnegative CHECK (height_cm >= 0);

-- +migrate Down
-- --- Rollback steps ---
ALTER TABLE observations DROP CONSTRAINT IF EXISTS height_nonnegative;
DROP INDEX IF EXISTS idx_observations_user_plant_date_desc;
DROP INDEX IF EXISTS idx_observations_plant_sort;
DROP INDEX IF EXISTS idx_plants_user_sort;
ALTER TABLE observations DROP COLUMN IF EXISTS sort_order;
ALTER TABLE plants DROP COLUMN IF EXISTS sort_order;
