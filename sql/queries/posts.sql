-- name: CreatePost :exec
-- Insert a new post into the `posts` table
INSERT INTO posts (
        id,
        -- Unique identifier for the post
        title,
        -- Title of the post
        url,
        -- URL of the post
        description,
        -- Brief description of the post
        published_at,
        -- Publication timestamp of the post
        feed_id -- Foreign key linking to the `feeds` table
    )
VALUES ($1, $2, $3, $4, $5, $6);
-- name: GetPostsForUser :many
-- Retrieve posts for all feeds followed by a specific user
SELECT posts.title,
    -- Title of the post
    posts.url,
    -- URL of the post
    posts.description,
    -- Description of the post
    posts.published_at -- Publication timestamp of the post
FROM posts
    INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1 -- Filter by the user ID
ORDER BY posts.published_at DESC -- Order posts by publication date, most recent first
LIMIT $2;
-- Limit the number of posts returned