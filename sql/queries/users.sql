-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
	$1,
	$2,
	$3,
	$4
)

RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ResetUsers :exec
DELETE FROM users WHERE id is not null;


-- name: GetUsers :many
SELECT * FROM users;
