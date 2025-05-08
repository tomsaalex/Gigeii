-- name: SelectUsers :many
select * from users;

-- name: AddUser :one
insert into users
(email, username, pass_hash, pass_salt, role)
values
($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserById :one
select * from users where id = $1;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: DeleteUser :exec
delete from users where id = $1;