
-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)

RETURNING *;

-- name: GetFeeds :many
SELECT f.name as feed_name, url, u.name as user_name FROM Feeds f JOIN users u ON f.user_id = u.id;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = $2, updated_at = $2 WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds order by last_fetched_at ASC NULLS FIRST LIMIT 1;
 


