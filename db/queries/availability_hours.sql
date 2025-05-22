-- name: AddAvailabilityHour :one
INSERT INTO availability_hours (availability_id, hour)
VALUES ($1, $2)
RETURNING *;