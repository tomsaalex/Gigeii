-- name: CreateAvailability :one
INSERT INTO availabilities (
   start_date, end_date,
  days, price, max_participants, precedance, created_by, duration
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;


-- name: GetAvailabilityByID :many
SELECT 
 a.id AS availability_id,
 start_date,
 end_date,
 days,
 price,
 max_participants,
 precedance,
 created_by,
 created_at,
 updated_at,
 duration,
 ah.id AS hour_id,
 hour
FROM availabilities a
INNER JOIN availability_hours ah ON ah.availability_id = a.id
WHERE a.id = $1;

-- name: UpdateAvailability :one
UPDATE availabilities SET 
  start_date = $2, 
  end_date = $3, 
  days = $4, 
  price = $5, 
  max_participants = $6, 
  precedance = $7, 
  created_by = $8, 
  duration = $9,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAvailability :one
DELETE FROM availabilities WHERE id = $1
RETURNING *;