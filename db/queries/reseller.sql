

-- name: AddReseller :one
INSERT INTO resellers (
    name, username, password_hash, email
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetResellerByID :one
SELECT * FROM resellers WHERE id = $1;

-- name: GetResellerByUsername :one
SELECT * FROM resellers WHERE username = $1;

-- name: GetResellerByEmail :one
SELECT * FROM resellers WHERE email = $1;


-- name: DeleteReseller :exec
DELETE FROM resellers WHERE id = $1;


-- name: SelectResellers :many
SELECT * FROM resellers;

