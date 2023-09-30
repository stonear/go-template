-- name: Index :many
SELECT id, name
FROM person
ORDER BY name;

-- name: Show :one
SELECT id, name
FROM person
WHERE id = $1
LIMIT 1;

-- name: Store :one
INSERT INTO person(name)
VALUES ($1)
RETURNING *;

-- name: Update :one
UPDATE person
SET name = $1
WHERE id = $2
RETURNING *;

-- name: Destroy :exec
DELETE FROM person
WHERE id = $1;