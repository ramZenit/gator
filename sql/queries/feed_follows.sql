-- name: CreateFeedFollow :one
WITH new_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT nff.*,
    f.name as feed_name,
    u.name as user_name
FROM new_feed_follow nff
    JOIN feeds f ON nff.feed_id = f.id
    JOIN users u ON nff.user_id = u.id;
-- name: GetFeedFollowsForUser :many
SELECT f.id,
    f.name as feed_name,
    f.url,
    u.name as user_name
FROM users u
    JOIN feed_follows ff ON u.id = ff.user_id
    JOIN feeds f ON ff.feed_id = f.id
WHERE u.name = $1;
-- name: DeleteFeedFollow :exec
WITH feed_user AS (
    SELECT u.id AS user_id,
        f.id AS feed_id
    FROM users u,
        feeds f
    WHERE u.name = $1
        AND f.url = $2
)
DELETE FROM feed_follows ff USING feed_user fu
WHERE ff.user_id = fu.user_id
    AND ff.feed_id = fu.feed_id;