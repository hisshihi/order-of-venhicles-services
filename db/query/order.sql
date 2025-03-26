-- name: CreateOrder :one
INSERT INTO "orders" (
        client_id,
        category_id,
        subtitle_category_id,
        status,
        client_message,
        order_date
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOrderByID :one
SELECT o.*,
    sc.name as category_name,
    suc.name as subtitle_category,
    u.username as client_name,
    u.phone as client_phone,
    u.whatsapp as client_whatsapp,
    u.city as client_city,
    u.district as client_district,
    p.username as provider_name,
    p.phone as provider_phone,
    p.whatsapp as provider_whatsapp,
    s.title as service_title
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
    JOIN "subtitle_category" suc ON o.subtitle_category_id = suc.id
    JOIN "users" u ON o.client_id = u.id
    LEFT JOIN "services" s ON s.id = o.service_id
    LEFT JOIN "users" p ON s.provider_id = p.id
WHERE o.id = $1;

-- name: ListOrdersByClientID :many
SELECT o.*,
    sc.name as category_name,
    suc.name as subcategory_name,
    s.title as service_title,
    p.username as provider_name
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
    JOIN "subtitle_category" suc ON o.subtitle_category_id = suc.id
    LEFT JOIN "services" s ON o.service_id = s.id
    LEFT JOIN "users" p ON s.provider_id = p.id
WHERE o.client_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCountOrdersByClientID :one
SELECT COUNT(*) FROM "orders"
WHERE client_id = $1;

-- name: ListAvailableOrdersForProvider :many
-- Получает список доступных заказов для провайдера услуг
SELECT o.*,
    sc.name as category_name,
    suc.name as subtitle_category,
    u.username as client_name,
    u.city as client_city,
    u.district as client_district
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
    JOIN "subtitle_category" suc ON o.subtitle_category_id = suc.id
    JOIN "users" u ON o.client_id = u.id
WHERE -- Заказ все еще открыт (pending)
    o.status = 'pending'
    -- AND -- Провайдер предлагает услуги в этой категории
    -- o.category_id IN (
    --     SELECT DISTINCT category_id
    --     FROM "services"
    --     WHERE provider_id = $1
    -- )
    AND -- Заказ не был принят провайдером
    o.provider_accepted = false
ORDER BY o.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListCountAvailableOrdersForProvider :one
SELECT COUNT(*) FROM "orders"
WHERE status = 'pending'
    AND provider_accepted = false
    AND category_id IN (
        SELECT DISTINCT category_id
        FROM "services"
        WHERE provider_id = $1
    );

-- name: AcceptOrderByProviderID :one
-- Провайдер принимает заказ и указывает свою услугу
UPDATE "orders"
SET provider_accepted = true,
    provider_message = $4,
    service_id = $3,
    status = 'accepted',
    updated_at = now()
WHERE id = $1
    AND -- Проверяем, что провайдер предлагает услуги в категории заказа
    (
        SELECT category_id
        FROM "orders"
        WHERE id = $1
    ) IN (
        SELECT DISTINCT category_id
        FROM "services"
        WHERE provider_id = $2
    )
    AND -- Заказ не должен быть уже принятым
    (
        provider_accepted = false
        OR provider_accepted IS NULL
    )
    AND -- Заказ должен иметь статус pending
    status = 'pending'
    AND -- Город сервиса должен совпадать с городом клиента или один из них должен быть NULL
    EXISTS (
        SELECT 1
        FROM "orders" o
            JOIN "users" client ON o.client_id = client.id
            JOIN "services" s ON s.id = $3
        WHERE o.id = $1
            AND s.provider_id = $2
            AND (
                -- Или город сервиса совпадает с городом клиента
                (
                    s.city IS NOT NULL
                    AND client.city IS NOT NULL
                    AND s.city = client.city
                ) -- Или сервис не привязан к городу (работает везде)
                OR s.city IS NULL -- Или клиент не указал город (не важно где)
                OR client.city IS NULL
            )
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
WITH provider_orders AS (
    SELECT o.*
    FROM "orders" o
    WHERE o.service_id IN (
            SELECT id
            FROM "services"
            WHERE provider_id = $1
        )
        OR (
            o.provider_accepted = true
            AND EXISTS (
                SELECT 1
                FROM "services" s
                WHERE s.provider_id = $1
                    AND s.category_id = o.category_id
            )
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
SET category_id = $2,
    client_message = $3,
    subtitle_category_id = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListOrders :many
SELECT * FROM "orders"
LIMIT $1 OFFSET $2;

-- name: CountOrders :one
SELECT COUNT(*) FROM "orders";

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: GetOrdersByCategory :many
-- Получает список заказов по категории
SELECT o.*,
    sc.name as category_name,
    suc.name as subtitle_category,
    u.username as client_name
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
    JOIN "subtitle_category" suc ON o.subtitle_category_id = suc.id
    JOIN "users" u ON o.client_id = u.id
WHERE o.category_id = $1
    AND o.status = 'pending'
    AND o.provider_accepted = false
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrdersBySubCategory :many
-- Получает список заказов по подкатегориям
SELECT o.*,
    sc.name as category_name,
    suc.name as subtitle_category,
    u.username as client_name
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
    JOIN "subtitle_category" suc ON o.subtitle_category_id = suc.id
    JOIN "users" u ON o.client_id = u.id
WHERE o.subtitle_category_id = $1
    AND o.status = 'pending'
    AND o.provider_accepted = false
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteOrdersByCategoryID :execrows
DELETE FROM orders
WHERE category_id = $1;

-- name: DeleteOrdersBySubcategoryID :execrows
DELETE FROM orders
WHERE subtitle_category_id = $1;