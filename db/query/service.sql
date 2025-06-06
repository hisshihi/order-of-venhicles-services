-- name: CreateService :one
INSERT INTO "services" (
        provider_id,
        category_id,
        subtitle_category_id,
        subcategory,
        title,
        description,
        price,
        country,
        city,
        district
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;
-- name: GetServiceByID :one
SELECT s.*,
    u.username as provider_name,
    u.photo_url as provider_photo,
    u.phone as provider_phone,
    u.whatsapp as provider_whatsapp,
    sc.name as category_name,
    (
        SELECT COUNT(*)
        FROM "reviews" r
        WHERE r.provider_id = s.provider_id
    ) as reviews_count,
    (
        SELECT COALESCE(AVG(r.rating), 0)
        FROM "reviews" r
        WHERE r.provider_id = s.provider_id
    ) as average_rating
FROM "services" s
    JOIN "users" u ON s.provider_id = u.id
    JOIN "service_categories" sc ON s.category_id = sc.id
WHERE s.id = $1
    AND u.is_blocked = false;
-- name: GetServicesByProviderID :many
SELECT s.*,
    sc.name as category_name
FROM "services" s
    JOIN "service_categories" sc ON s.category_id = sc.id
    JOIN "users" u ON s.provider_id = u.id
WHERE s.provider_id = $1
    AND u.is_blocked = false
ORDER BY s.created_at DESC
LIMIT $2 OFFSET $3;
-- name: UpdateService :one
UPDATE "services"
SET category_id = $3,
    subcategory = $4,
    title = $5,
    description = $6,
    price = $7,
    country = $8,
    city = $9,
    district = $10,
    subtitle_category_id = $11,
    updated_at = NOW()
WHERE id = $1
    AND provider_id = $2
RETURNING *;
-- name: DeleteService :exec
DELETE FROM "services"
WHERE id = $1
    AND provider_id = $2;
-- name: ListServices :many
SELECT s.*,
    u.username as provider_name,
    u.photo_url as provider_photo,
    sc.name as category_name
FROM "services" s
    JOIN "users" u ON s.provider_id = u.id
    JOIN "service_categories" sc ON s.category_id = sc.id
WHERE u.is_blocked = false
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;
-- name: ListCountServicesByCatetegory :one
SELECT COUNT(*)
FROM "services"
WHERE category_id = $1;
-- name: ListServicesByCategory :many
SELECT s.*,
    u.username as provider_name,
    u.photo_url as provider_photo,
    sc.name as category_name
FROM "services" s
    JOIN "users" u ON s.provider_id = u.id
    JOIN "service_categories" sc ON s.category_id = sc.id
WHERE s.category_id = $1
    AND u.is_blocked = false
ORDER BY s.created_at DESC
LIMIT $2 OFFSET $3;
-- name: ListServicesByLocation :many
SELECT s.*,
    u.username as provider_name,
    u.photo_url as provider_photo,
    sc.name as category_name
FROM "services" s
    JOIN "users" u ON s.provider_id = u.id
    JOIN "service_categories" sc ON s.category_id = sc.id
WHERE (
        $1::text IS NULL
        OR s.country = $1
    )
    AND (
        $2::text IS NULL
        OR s.city = $2
    )
    AND (
        $3::text IS NULL
        OR s.district = $3
    )
    AND u.is_blocked = false
ORDER BY s.created_at DESC
LIMIT $4 OFFSET $5;
-- name: SearchServices :many
SELECT s.*,
    u.username as provider_name,
    u.photo_url as provider_photo,
    sc.name as category_name
FROM "services" s
    JOIN "users" u ON s.provider_id = u.id
    JOIN "service_categories" sc ON s.category_id = sc.id
WHERE (
        to_tsvector('simple', s.title) @@ to_tsquery('simple', $1)
        OR to_tsvector('simple', s.description) @@ to_tsquery('simple', $1)
        OR s.title ILIKE '%' || $1 || '%'
        OR s.description ILIKE '%' || $1 || '%'
    )
    AND u.is_blocked = false
ORDER BY s.created_at DESC
LIMIT $2 OFFSET $3;
-- name: ListServicesByProviderIDAndCategory :many
SELECT s.*,
    sc.name as category_name
FROM "services" s
    JOIN "service_categories" sc ON s.category_id = sc.id
    JOIN "users" u ON s.provider_id = u.id
WHERE s.provider_id = $1
    AND s.category_id = $2
    AND u.is_blocked = false;
-- name: ListServicesByProviderIDAndSubCategory :many
SELECT s.*,
    sc.name as subcategory_name
FROM "services" s
    JOIN "subtitle_category" sc ON s.subtitle_category_id = sc.id
    JOIN "users" u ON s.provider_id = u.id
WHERE s.provider_id = $1
    AND s.subtitle_category_id = $2
    AND u.is_blocked = false;
-- name: CountServicesByProviderID :one
SELECT COUNT(*)
FROM "services"
WHERE provider_id = $1;
-- Фильтрация услуг по цене
-- name: FilterServiceByPrice :many
SELECT *
FROM "services"
WHERE price >= $1
    AND price <= $2;
-- name: CountService :one
SELECT COUNT(*)
FROM "services";
-- name: DeleteServicesByCategoryID :execrows
DELETE FROM services
WHERE category_id = $1;
-- name: DeleteServicesBySubcategoryID :execrows
DELETE FROM services
WHERE subtitle_category_id = $1;