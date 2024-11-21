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