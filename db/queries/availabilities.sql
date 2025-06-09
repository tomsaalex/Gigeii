-- name: CreateAvailability :one
INSERT INTO availabilities (
   start_date, end_date,
  days, price, max_participants, precedance, created_by, duration, notes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: FindAvailabilityConflicts :many
SELECT a.id, a.start_date, a.end_date, a.days, a.price, a.max_participants, a.precedance, a.created_by, a.created_at, a.updated_at, a.duration, ah.hour 
FROM availabilities a
JOIN availability_hours ah ON ah.availability_id = a.id
WHERE 
  a.start_date <= @start_date
  AND a.end_date >= @end_date
  AND (a.days & @days) != 0 
  AND ah.hour = ANY(@hours::timestamptz[])
  AND (@availability_id::uuid IS NULL OR a.id != @availability_id::uuid)
;

-- name: ShiftPrecedenceAbove :exec
UPDATE availabilities SET
precedance = precedance + 1
WHERE precedance > $1; 

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
 notes,
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
  notes = $10,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAvailability :one
DELETE FROM availabilities WHERE id = $1
RETURNING *;

-- name: GetAllAvailabilities :many
SELECT 
    a.id,
    a.start_date,
    a.end_date,
    a.days,
    a.price,
    a.max_participants,
    a.precedance,
    a.created_by,
    a.created_at,
    a.updated_at,
    a.duration,
    a.notes,
    ah.hour
FROM availabilities a
LEFT JOIN availability_hours ah ON ah.availability_id = a.id;