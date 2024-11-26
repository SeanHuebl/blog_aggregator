-- +goose Up
-- Create the `users` table to store user details
CREATE TABLE users (
    id UUID PRIMARY KEY,
    -- Unique identifier for the user
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp for when the user record was created
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp for the last update to the user record
    name TEXT UNIQUE NOT NULL -- Unique username for the user
);
-- +goose Down
-- Drop the `users` table and all associated data
DROP TABLE users CASCADE;