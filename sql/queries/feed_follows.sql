-- name: CreateFeedFollow :many
-- Create a new feed follow relationship and return the details
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, user_id, feed_id) -- Insert a new follow relationship
    VALUES ($1, $2, $3)
    RETURNING * -- Return the inserted follow record
)
SELECT inserted_feed_follow.*,
    -- Include all fields from the inserted follow record
    feeds.name AS feed_name,
    -- Include the name of the feed being followed
    users.name AS user_name -- Include the name of the user following the feed
FROM inserted_feed_follow
    INNER JOIN users ON users.id = inserted_feed_follow.user_id
    INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;
-- name: GetFeedFollowsForUser :many
-- Retrieve all feeds followed by a specific user with detailed information
SELECT *,
    -- Include all fields from the feed_follows table
    feeds.name AS feed_name,
    -- Include the name of the feed being followed
    users.name AS user_name -- Include the name of the user following the feed
FROM feed_follows
    INNER JOIN users ON users.id = feed_follows.user_id
    INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE users.id = $1;
-- Filter by the user ID
-- name: Unfollow :exec
-- Remove a feed follow relationship for a specific user and feed
DELETE FROM feed_follows
WHERE user_id = $1 -- Specify the user ID
    AND feed_id = $2;
-- Specify the feed ID