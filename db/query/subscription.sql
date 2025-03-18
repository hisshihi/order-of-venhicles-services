-- name: CreateSubscription :one
INSERT INTO subscriptions (provider_id, start_date, subscription_type, promo_code_id, price, end_date, status, original_price)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: GetSubscriptionByID :one
SELECT *
FROM subscriptions
WHERE id = $1;
-- name: GetSubscriptionByProviderID :one
SELECT *
FROM subscriptions
WHERE provider_id = $1;
-- name: ListSubscriptions :many
SELECT *
FROM subscriptions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
-- name: ListSubscriptionsByProviderID :many
SELECT *
FROM subscriptions
WHERE provider_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
-- name: UpdateSubscription :one
UPDATE subscriptions
SET provider_id = $2,
    start_date = $3,
    end_date = $4,
    status = $5,
    subscription_type = $6,
    price = $7,
    original_price = $8,
    promo_code_id = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE id = $1;
-- name: GetActiveSubscriptionForProvider :one
SELECT *
FROM subscriptions
WHERE provider_id = $1
    AND status = 'active'
    AND end_date > NOW()
ORDER BY end_date DESC
LIMIT 1;
-- name: CheckAndUpdateExpiredSubscriptions :many
UPDATE subscriptions
SET status = 'expired',
    updated_at = NOW()
WHERE status = 'active'
    AND end_date <= NOW()
RETURNING *;