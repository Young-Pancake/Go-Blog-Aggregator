-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;


-- name: GetUsers :many
SELECT * FROM users;

-- name: DeletAllUsers :exec
DELETE FROM users;