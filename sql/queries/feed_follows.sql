-- name: CreateFeedFollow :one
WITH new_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
)
SELECT
    nff.*,
    f.name as feed_name,
    u.name as user_name
FROM new_feed_follow nff
JOIN feeds f ON nff.feed_id = f.id
JOIN users u ON nff.user_id = u.id;


-- name: GetFeedFollowsForUser :many
SELECT f.id, f.name as feed_name, f.url, u.name as user_name 
FROM users u
JOIN feed_follows ff ON u.id = ff.user_id 
JOIN feeds f ON ff.feed_id = f.id
WHERE u.name = $1; 