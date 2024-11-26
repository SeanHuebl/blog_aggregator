-- name: CreateUser :one
-- Insert a new user into the `users` table
-- Returns the created user record
INSERT INTO users (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetUser :one
-- Retrieve a user by their username
SELECT *
FROM users
WHERE name = $1;
-- name: GetUsers :many
-- Retrieve the list of all usernames
SELECT name
FROM users;
-- name: Reset :exec
-- Delete all user records from the `users` table
DELETE FROM users;