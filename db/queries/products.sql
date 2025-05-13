-- name: CreateProduct :one
INSERT INTO products (
    name,
    timezone,
    created_by
)
VALUES (
    $1, $2, $3
)
RETURNING id, name, timezone;
