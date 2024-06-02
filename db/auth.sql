-- name: Show :one
SELECT id, username, password
FROM auth
WHERE username = $1
LIMIT 1;
