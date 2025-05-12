-- name: CreateAvailability :one
INSERT INTO availabilities (
  product_id, start_date, end_date, availability_type,
  days, hours, price, max_participants, precedance, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;
