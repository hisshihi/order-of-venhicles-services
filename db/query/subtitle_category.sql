-- name: CreateSubtitle :one
INSERT INTO subtitle_category (name)
VALUES ($1)
RETURNING *;

-- name: GetSubtitleCategoryByID :one
SELECT * FROM subtitle_category
WHERE id = $1;

-- name: ListSubtitleCategory :many
SELECT * FROM subtitle_category
ORDER BY name ASC;

-- name: UpdateSubtitleCategory :one
UPDATE subtitle_category
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteSubtitleCategory :execrows
DELETE FROM subtitle_category
WHERE id = $1;