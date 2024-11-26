-- +goose Up
-- Create the `posts` table to store posts fetched from RSS feeds
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    -- Unique identifier for the post
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Creation timestamp
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Last update timestamp
    title TEXT,
    -- Title of the post (optional)
    url TEXT UNIQUE,
    -- Unique URL of the post
    description TEXT,
    -- Summary or description of the post (optional)
    published_at TIMESTAMP,
    -- Publication timestamp of the post
    feed_id UUID NOT NULL,
    -- Foreign key linking to the `feeds` table
    CONSTRAINT feed_fk FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE -- Cascade delete on feed removal
);
-- +goose Down
-- Drop the `posts` table and all its associated data
DROP TABLE posts CASCADE;