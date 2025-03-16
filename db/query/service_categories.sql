-- name: CreateServiceCategory :one
INSERT INTO service_categories (name)
VALUES ($1)
RETURNING *;

-- name: GetServiceCategoryByID :one
SELECT * FROM service_categories
WHERE id = $1;

-- name: ListServiceCategories :many
SELECT * FROM service_categories
ORDER BY id ASC
LIMIT $1 OFFSET $2;

-- name: UpdateServiceCategory :one
UPDATE service_categories
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteServiceCategory :exec
DELETE FROM service_categories
WHERE id = $1;