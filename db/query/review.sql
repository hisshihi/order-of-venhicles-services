-- name: CreateReview :one
INSERT INTO reviews (order_id, client_id, provider_id, rating, comment)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetReviewByID :one
SELECT * FROM reviews
WHERE id = $1;

-- name: ListReviewsByProviderID :many
SELECT * FROM reviews
WHERE provider_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProviderOverallRating :one
SELECT provider_id, AVG(rating) AS overall_rating, COUNT(*) AS review_count
FROM reviews
WHERE provider_id = $1
GROUP BY provider_id;

-- name: UpdateReview :one
UPDATE reviews
SET rating = $2, comment = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;