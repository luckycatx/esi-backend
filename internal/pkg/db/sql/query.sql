-- name: CreateUser :exec
INSERT INTO users (username, email, pwd)
VALUES (?, ?, ?);
-- name: FetchUsers :many
SELECT *
FROM users
ORDER BY id;
-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;
-- name: UpdateUser :exec
UPDATE users
SET username = ?,
  email = ?,
  pwd = ?
WHERE id = ?