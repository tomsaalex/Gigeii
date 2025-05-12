-- +goose Up
CREATE TABLE IF NOT EXISTS availabilities (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    -- Period definition
    start_date DATE NOT NULL,
    end_date DATE,

    -- Type: fixed = exact dates, daily = every day, weekly = selected days
    availability_type TEXT NOT NULL CHECK (availability_type IN ('fixed', 'daily', 'weekly')),

    days INTEGER NOT NULL CHECK (days >= 0 AND days <= 127),            -- 7-bit mask (0=Sun to 6=Sat)
    hours INTEGER NOT NULL CHECK (hours >= 0 AND hours <= 16777215),    -- 24-bit mask (0=00:00 to 23:00)
    
    -- Default config for all generated openings
    price INTEGER NOT NULL CHECK (price >= 0),
    max_participants INTEGER NOT NULL CHECK (max_participants > 0),

    -- Conflict resolution
    precedance INTEGER NOT NULL DEFAULT 0,

    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS availabilities;
