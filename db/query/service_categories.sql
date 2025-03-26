-- name: CreateServiceCategory :one
INSERT INTO service_categories (name, icon, description, slug)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetServiceCategoryByID :one
SELECT * FROM service_categories
WHERE id = $1;

-- name: GetServiceCategoryBySlug :one
SELECT * FROM service_categories
WHERE slug = $1;

-- name: ListServiceCategories :many
SELECT * FROM service_categories
ORDER BY name ASC;

-- name: UpdateServiceCategory :one
UPDATE service_categories
SET name = $2, icon = $3, description = $4, slug = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteServiceCategory :execrows
DELETE FROM service_categories
WHERE id = $1;