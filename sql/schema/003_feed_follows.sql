-- +goose Up
-- Create the `feed_follows` table to track user subscriptions to feeds
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    -- Unique identifier for the follow relationship
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp for when the relationship was created
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Last updated timestamp
    user_id UUID NOT NULL,
    -- Foreign key linking to the `users` table
    feed_id UUID NOT NULL,
    -- Foreign key linking to the `feeds` table
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    -- Cascade delete on user removal
    CONSTRAINT feed_fk FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
    -- Cascade delete on feed removal
    UNIQUE (user_id, feed_id) -- Ensure each user can follow a feed only once
);
-- +goose Down
-- Drop the `feed_follows` table and all associated data
DROP TABLE feed_follows CASCADE;