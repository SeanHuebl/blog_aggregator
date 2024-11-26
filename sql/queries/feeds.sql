-- name: AddFeed :one
-- Insert a new feed into the `feeds` table
-- Returns the newly created feed record
INSERT INTO feeds (id, name, url, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetFeeds :many
-- Retrieve all feeds with their associated user names
SELECT feeds.name,
    -- Name of the feed
    feeds.url,
    -- URL of the feed
    users.name -- Name of the user who added the feed
FROM feeds
    INNER JOIN users ON feeds.user_id = users.id;
-- name: GetFeed :one
-- Retrieve a feed by its URL
SELECT id,
    -- Unique identifier for the feed
    name -- Name of the feed
FROM feeds
WHERE url = $1;
-- name: MarkFeedFetched :exec
-- Update the `last_fetched_at` and `updated_at` timestamps for a feed
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
    -- Set the last fetched timestamp to now
    updated_at = CURRENT_TIMESTAMP -- Update the modified timestamp to now
WHERE id = $1;
-- name: GetNextFeedToFetch :one
-- Retrieve the next feed to fetch, prioritizing the least recently fetched
SELECT id,
    -- Unique identifier for the feed
    url -- URL of the feed
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST -- Order by least recently fetched (NULLs come first)
LIMIT 1;