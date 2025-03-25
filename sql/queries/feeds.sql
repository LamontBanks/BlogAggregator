-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

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
SELECT feeds.name AS feed_name, feeds.url, users.name AS user_name
FROM feeds
LEFT JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.name ASC;