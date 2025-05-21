-- name: CreateAvailability :one
INSERT INTO availabilities (
   start_date, end_date,
  days, price, max_participants, precedance, created_by, duration
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;


