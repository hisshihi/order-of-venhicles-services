-- name: CreatePayment :one
INSERT INTO payments (user_id, amount, payment_method, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePayment :one
UPDATE payments
SET user_id = $2, amount = $3, payment_method = $4, status = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;