-- +goose Up
-- SQL section to apply the migration

-- 1. Drop the product_id column if it exists
ALTER TABLE availabilities DROP COLUMN IF EXISTS product_id;

-- 2. Drop the products table
DROP TABLE IF EXISTS products;

-- 3. Drop availability_type column
ALTER TABLE availabilities DROP COLUMN IF EXISTS availability_type;

-- 4. Drop hours column
ALTER TABLE availabilities DROP COLUMN IF EXISTS hours;

-- 5. Create availability_hours table
CREATE TABLE IF NOT EXISTS availability_hours (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    availability_id UUID NOT NULL REFERENCES availabilities(id) ON DELETE CASCADE,
    hour TIMESTAMPTZ NOT NULL,
    UNIQUE (availability_id, hour) -- Ensure no duplicate hours per availability
);

ALTER TABLE availabilities ADD COLUMN duration INTERVAL;

ALTER TABLE availabilities ADD COLUMN notes TEXT;


-- +goose Down
-- SQL section to rollback the migration

-- 6. Drop availability_hours table
DROP TABLE IF EXISTS availability_hours;

ALTER TABLE availabilities DROP COLUMN IF EXISTS duration;


ALTER TABLE availabilities DROP COLUMN IF EXISTS notes;