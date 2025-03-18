-- name: CreatePendingSubscription :one
INSERT INTO pending_subscriptions (
        payment_id,
        user_id,
        subscription_type,
        start_date,
        end_date,
        original_price,
        final_price,
        promo_code_id,
        is_update
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9
    )
RETURNING *;
-- name: GetPendingSubscriptionByPaymentID :one
SELECT *
FROM pending_subscriptions
WHERE payment_id = $1;
-- name: DeletePendingSubscriptionByPaymentID :exec
DELETE FROM pending_subscriptions
WHERE payment_id = $1;