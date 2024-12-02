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
SELECT * FROM users
WHERE users.id = (
  SELECT user_id FROM refresh_tokens
  WHERE refresh_tokens.id = $1
);


-- name: SetRevokedAt :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE id = $1;

