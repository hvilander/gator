-- name: CreateFeedFollow :one 
WITH new_feed_follow AS (
	INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
	VALUES ( $1, $2, $3, $4, $5)
	RETURNING *
)

SELECT new_feed_follow.*, f.name AS feed_name, u.name as user_name
FROM new_feed_follow
INNER JOIN users u ON new_feed_follow.user_id = u.id
INNER JOIN feeds f on new_feed_follow.feed_id = f.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFollowsByUserID :many
SELECT ff.*, u.name AS user_name, f.name AS feed_name
FROM feed_follows ff
INNER JOIN users u on ff.user_id = u.id
INNER JOIN feeds f on ff.feed_id = f.id
WHERE ff.user_id = $1; 


-- name: DeleteFeedFollowByUserAndURL :exec
DELETE FROM feed_follows ff
WHERE ff.user_id = $1 and ff.feed_id = (
	Select id from feeds where url = $2
);

