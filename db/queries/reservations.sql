-- name: ReserveOrUpdateReservation :one
INSERT INTO reservations (
    reservation_reference, external_reservation_reference, reseller_id, availability_id, date_time, quantity, status
)
VALUES (
    $1, $2, $3, $4, $5, $6, 'CONFIRMED'
)
ON CONFLICT (reseller_id, external_reservation_reference) DO UPDATE 
SET 
    reservation_reference = EXCLUDED.reservation_reference,
    availability_id = EXCLUDED.availability_id,
    date_time = EXCLUDED.date_time,
    quantity = EXCLUDED.quantity,
    status = 'CONFIRMED',
    updated_at = now()
RETURNING id, reservation_reference, external_reservation_reference, reseller_id, availability_id, date_time, quantity, status, created_at, updated_at;

-- name: GetReservationByExternalReference :one
SELECT id, reservation_reference, external_reservation_reference, reseller_id, availability_id, date_time, quantity, status, created_at, updated_at
FROM reservations
WHERE reseller_id = $1 AND external_reservation_reference = $2;

-- name: GetReservationByReference :one
SELECT id, reservation_reference, external_reservation_reference, reseller_id, availability_id, date_time, quantity, status, created_at, updated_at
FROM reservations
WHERE reservation_reference = $1;

-- name: CancelReservation :one
UPDATE reservations
SET status = 'CANCELED', updated_at = now()
WHERE reservation_reference = $1
RETURNING id, reservation_reference, external_reservation_reference, reseller_id, availability_id, date_time, quantity, status, created_at, updated_at;


