-- name: CreateSubscription :one
INSERT INTO subscriptions (provider_id, start_date, end_date, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSubscriptionByID :one
SELECT * FROM subscriptions
WHERE id = $1;

-- name: ListSubscriptions :many
SELECT * FROM subscriptions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET provider_id = $2, start_date = $3, end_date = $4, status = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE id = $1;

