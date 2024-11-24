-- name: AddFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetFeeds :many
SELECT feeds.name,
    feeds.url,
    users.name
FROM feeds
    INNER JOIN users ON feeds.user_id = users.id;
-- name: GetFeed :one
SELECT id,
    name
FROM feeds
WHERE url = $1;
-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
-- name: GetNextFeedToFetch :one
SELECT id
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;