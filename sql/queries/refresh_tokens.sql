-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(
  id,
  created_at,
  updated_at,
  user_id,
  expires_at
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
WHERE id = $1;


-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.id = $1
AND revoked_at is NULL
AND expires_at > NOW();


-- name: SetRevokedAt :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

