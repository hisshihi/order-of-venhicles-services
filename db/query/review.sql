-- name: CreateReview :one
-- Создает новый отзыв от клиента об услугодателе
INSERT INTO "reviews" (
        client_id,
        provider_id,
        rating,
        comment
    )
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetReviewsByProviderID :many
-- Получает все отзывы об услугодателе
SELECT r.id,
    r.client_id,
    u.username as client_name,
    u.photo_url as client_photo,
    r.rating,
    r.comment,
    r.created_at
FROM "reviews" r
    JOIN "users" u ON r.client_id = u.id
WHERE r.provider_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetReviewByID :one
-- Получает конкретный отзыв по ID
SELECT r.id,
    r.client_id,
    r.provider_id,
    r.rating,
    r.comment,
    r.created_at,
    uc.username as client_name,
    up.username as provider_name
FROM "reviews" r
    JOIN "users" uc ON r.client_id = uc.id
    JOIN "users" up ON r.provider_id = up.id
WHERE r.id = $1;

-- name: GetAverageRatingForProvider :one
-- Получает среднюю оценку услугодателя
SELECT COALESCE(AVG(rating), 0) as average_rating,
    COUNT(*) as total_reviews
FROM "reviews"
WHERE provider_id = $1;

-- name: DeleteReview :exec
-- Удаляет отзыв (только если пользователь является автором или администратором)
DELETE FROM "reviews"
WHERE id = $1
    AND client_id = $2;

-- name: CheckIfClientReviewedOrder :one
-- Проверяет, оставил ли клиент отзыв по данному заказу
SELECT EXISTS(
        SELECT 1
        FROM "reviews"
        WHERE client_id = $1
            AND provider_id = $2
    ) as has_review;