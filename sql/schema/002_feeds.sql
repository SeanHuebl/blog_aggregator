-- +goose Up
-- Create the `feeds` table to store RSS feed details
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    -- Unique identifier for the feed
    name TEXT NOT NULL,
    -- Name of the feed (e.g., "Tech News")
    url TEXT UNIQUE NOT NULL,
    -- Unique URL of the feed
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp for when the feed was created
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp for the last update to the feed
    user_id UUID NOT NULL,
    -- Foreign key linking to the `users` table
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE -- Cascade delete on user removal
);
-- +goose Down
-- Drop the `feeds` table and all associated data
DROP TABLE feeds CASCADE;