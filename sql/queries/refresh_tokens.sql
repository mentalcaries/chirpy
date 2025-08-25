-- name: CreateRefreshToken :one
INSERT into refresh_tokens (
    token, created_at, updated_at, user_id, expires_at
)
VALUES (
    $1, 
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;