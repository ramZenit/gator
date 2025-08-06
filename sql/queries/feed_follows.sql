-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;


-- name: GetFeedFollowsForUser :many
SELECT f.id, f.name as feed_name, f.url, u.name as user_name 
FROM users u
JOIN feed_follows ff ON u.id = ff.user_id 
JOIN feeds f ON ff.feed_id = f.id
WHERE u.name = $1; 