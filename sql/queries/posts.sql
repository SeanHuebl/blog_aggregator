-- name: CreatePost :exec
INSERT INTO posts (
        id,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES ($1, $2, $3, $4, $5, $6);
-- name: GetPostsForUser :many
SELECT posts.title,
    posts.url,
    posts.description,
    posts.published_at
FROM posts
    INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;