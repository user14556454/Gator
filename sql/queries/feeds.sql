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
SELECT 
    name,
    url,
    user_id
FROM
    feeds;

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- name: GetFeed :one
SELECT 
    *
FROM
    feeds
WHERE
    url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
    updated_at = $1,
    last_fetched_at = $2
WHERE
    user_id = $3;

-- name: GetNextFeedToFetch :one
SELECT 
    *
FROM 
    feeds
ORDER BY
    last_fetched_at ASC NULLS FIRST
LIMIT
    1;