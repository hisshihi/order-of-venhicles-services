-- name: CreateService :one
INSERT INTO services (provider_id, category_id, title, description, price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetServiceByID :one
SELECT * FROM services
WHERE id = $1;

-- name: ListServicesByTitle :many
SELECT * FROM services
WHERE title ILIKE '%' || $1 || '%'
ORDER BY title;

-- name: ListServicesByCategory :many
SELECT * FROM services
WHERE category_id = $1
ORDER BY title;


-- name: ListServices :many
SELECT * FROM services
ORDER BY title DESC
LIMIT $1 OFFSET $2;

-- name: UpdateService :one
UPDATE services
SET provider_id = $2, category_id = $3, title = $4, description = $5, price = $6
WHERE id = $1
RETURNING *;

-- name: DeleteService :exec
DELETE FROM services
WHERE id = $1;