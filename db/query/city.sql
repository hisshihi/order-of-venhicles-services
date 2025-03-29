-- name: CreateCity :one
INSERT INTO city (name)
VALUES ($1)
RETURNING *;

-- name: ListCity :many
SELECT * FROM city
ORDER BY name;

-- name: GetCityByID :one
SELECT * FROM city
WHERE id = $1;

-- name: UpdateCity :one
UPDATE city
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCity :exec
DELETE FROM city
WHERE id = $1;