-- name: CreateOrder :one
INSERT INTO "orders" (
        client_id,
        service_id,
        status,
        client_message,
        order_date
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: GetOrderByID :one
SELECT o.*,
    s.title as service_title,
    u.username as client_name,
    u.phone as client_phone,
    u.whatsapp as client_whatsapp,
    p.username as provider_name,
    p.phone as provider_phone,
    p.whatsapp as provider_whatsapp
FROM "orders" o
    JOIN "services" s ON o.service_id = s.id
    JOIN "users" u ON o.client_id = u.id
    JOIN "users" p ON s.provider_id = p.id
WHERE o.id = $1;
-- name: ListOrdersByClientID :many
SELECT o.*,
    s.title as service_title,
    p.username as provider_name
FROM "orders" o
    JOIN "services" s ON o.service_id = s.id
    JOIN "users" p ON s.provider_id = p.id
WHERE o.client_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;
-- name: ListAvailableOrdersForProvider :many
-- Получает список доступных заказов для провайдера услуг
SELECT o.*,
    s.title as service_title,
    c.name as category_name,
    u.username as client_name
FROM "orders" o
    JOIN "services" s ON o.service_id = s.id
    JOIN "service_categories" c ON s.category_id = c.id
    JOIN "users" u ON o.client_id = u.id
WHERE -- Заказ все еще открыт (pending)
    o.status = 'pending'
    AND -- Провайдер предлагает услуги в этой категории
    s.category_id IN (
        SELECT DISTINCT category_id
        FROM "services"
        WHERE provider_id = $1
    )
    AND -- Заказ не был принят провайдером
    o.provider_accepted = false
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;
-- name: AcceptOrderByProvider :one
-- Провайдер принимает заказ
UPDATE "orders"
SET provider_accepted = true,
    provider_message = $3,
    status = 'accepted',
    updated_at = now()
WHERE id = $1
    AND -- Проверяем, что провайдер предлагает услуги в категории заказа
    (
        SELECT category_id
        FROM "services"
        WHERE id = (
                SELECT service_id
                FROM "orders"
                WHERE id = $1
            )
    ) IN (
        SELECT DISTINCT category_id
        FROM "services"
        WHERE provider_id = $2
    )
RETURNING *;
-- name: UpdateOrderStatus :one
UPDATE "orders"
SET status = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;
-- name: GetOrderStatistics :one
-- Получает статистику заказов для услугодателя
WITH provider_services AS (
    SELECT id
    FROM "services"
    WHERE provider_id = $1
),
provider_orders AS (
    SELECT *
    FROM "orders"
    WHERE service_id IN (
            SELECT id
            FROM provider_services
        )
)
SELECT COUNT(*) FILTER (
        WHERE status = 'pending'
    ) as pending_count,
    COUNT(*) FILTER (
        WHERE status = 'accepted'
    ) as accepted_count,
    COUNT(*) FILTER (
        WHERE status = 'completed'
    ) as completed_count,
    COUNT(*) FILTER (
        WHERE status = 'cancelled'
    ) as cancelled_count,
    COUNT(*) as total_count
FROM provider_orders;
-- name: UpdateOrder :one
UPDATE orders
SET client_id = $2,
    service_id = $3,
    status = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;