-- name: AddProviderToFavorites :one
-- Добавляет услугодателя в избранное клиента
INSERT INTO "favorites" (client_id, provider_id)
VALUES ($1, $2) ON CONFLICT (client_id, provider_id) DO NOTHING
RETURNING *;
-- name: RemoveProviderFromFavorites :exec
-- Удаляет услугодателя из избранного клиента
DELETE FROM "favorites"
WHERE client_id = $1
    AND provider_id = $2;
-- name: ListFavoriteProviders :many
-- Получает список избранных услугодателей клиента
SELECT f.id,
    f.provider_id,
    u.username as provider_name,
    u.photo_url as provider_photo,
    u.phone as provider_phone,
    u.whatsapp as provider_whatsapp,
    u.description as provider_description,
    (
        SELECT COUNT(*)
        FROM "reviews" r
        WHERE r.provider_id = f.provider_id
    ) as reviews_count,
    (
        SELECT COALESCE(AVG(r.rating), 0)
        FROM "reviews" r
        WHERE r.provider_id = f.provider_id
    ) as average_rating,
    f.created_at
FROM "favorites" f
    JOIN "users" u ON f.provider_id = u.id
WHERE f.client_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;
-- name: CheckIfProviderIsFavorite :one
-- Проверяет, добавлен ли услугодатель в избранное клиента
SELECT EXISTS(
        SELECT 1
        FROM "favorites"
        WHERE client_id = $1
            AND provider_id = $2
    ) as is_favorite;