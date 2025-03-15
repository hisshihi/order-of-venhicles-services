-- name: CreateReview :one
INSERT INTO reviews (order_id, client_id, provider_id, rating, comment)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetReviewByID :one
SELECT * FROM reviews
WHERE id = $1;

-- name: ListReviews :many
SELECT * FROM reviews
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: UpdateReview :one
UPDATE reviews
SET rating = $2, comment = $3
WHERE id = $1
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;