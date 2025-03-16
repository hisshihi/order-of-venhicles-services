-- name: CreateOrder :one
INSERT INTO orders (client_id, service_id, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET client_id = $2, service_id = $3, status = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;