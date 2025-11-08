-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM refresh_tokens
JOIN users
ON refresh_tokens.user_id = users.id
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET updated_at = $2,
    revoked_at = $3
WHERE token = $1;
