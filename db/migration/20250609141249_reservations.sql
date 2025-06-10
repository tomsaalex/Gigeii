-- +goose Up
CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    reservation_reference TEXT NOT NULL UNIQUE, -- ID intern 
    external_reservation_reference TEXT NOT NULL,
    reseller_id UUID NOT NULL REFERENCES resellers(id) ON DELETE CASCADE,

    date_time TIMESTAMPTZ NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),

    status TEXT NOT NULL CHECK (status IN ('CONFIRMED', 'CANCELED')),

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Asigură unicitatea unei referințe externe per reseller
CREATE UNIQUE INDEX IF NOT EXISTS uniq_reseller_ext_ref
ON reservations (reseller_id, external_reservation_reference);

-- +goose Down
DROP TABLE IF EXISTS reservations;
