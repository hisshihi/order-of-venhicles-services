-- name: CreateOrderResponse :one
-- Создает новый отклик на заказ от провайдера услуг
-- Входные параметры: ID заказа, ID провайдера, сообщение, предложенная цена, выбран ли отклик
-- Возвращает: созданную запись отклика
INSERT INTO order_responses (
    order_id, 
    provider_id, 
    message, 
    offered_price, 
    is_selected
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetOrderResponseById :one
-- Получает отклик по его ID
-- Входные параметры: ID отклика
-- Возвращает: запись отклика или ничего, если отклик не найден
SELECT * FROM order_responses
WHERE id = $1
LIMIT 1;

-- name: GetOrderResponsesByOrderId :many
-- Получает все отклики для конкретного заказа с информацией о провайдерах
-- Входные параметры: ID заказа
-- Возвращает: список откликов с данными провайдеров (имя, телефон, whatsapp, фото)
SELECT r.*, u.username, u.phone, u.whatsapp, u.photo_url
FROM order_responses r
JOIN users u ON r.provider_id = u.id
WHERE r.order_id = $1
ORDER BY r.created_at;

-- name: GetOrderResponsesByProviderId :many
-- Получает все отклики, созданные конкретным провайдером, с информацией о связанных заказах
-- Входные параметры: ID провайдера, лимит и смещение для пагинации
-- Возвращает: список откликов с основной информацией о заказах
SELECT r.*, o.client_id, o.category_id, o.service_id, o.status, o.client_message, o.order_date
FROM order_responses r
JOIN orders o ON r.order_id = o.id
WHERE r.provider_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrderResponseByOrderAndProvider :one
-- Проверяет наличие отклика от конкретного провайдера на конкретный заказ
-- Входные параметры: ID заказа, ID провайдера
-- Возвращает: запись отклика или ничего, если отклик не найден
SELECT * FROM order_responses
WHERE order_id = $1 AND provider_id = $2
LIMIT 1;

-- name: CountOrderResponsesByOrderId :one
-- Подсчитывает количество откликов на конкретный заказ
-- Входные параметры: ID заказа
-- Возвращает: количество откликов (число)
SELECT COUNT(*) FROM order_responses
WHERE order_id = $1;

-- name: UpdateOrderResponse :one
-- Обновляет информацию в существующем отклике
-- Входные параметры: ID отклика, новое сообщение, новая предложенная цена
-- Возвращает: обновленную запись отклика
UPDATE order_responses
SET 
    message = $2,
    offered_price = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteOrderResponse :exec
-- Удаляет отклик по его ID
-- Входные параметры: ID отклика
-- Возвращает: ничего
DELETE FROM order_responses
WHERE id = $1;

-- name: SelectProviderForOrder :one
-- Выбирает провайдера для выполнения заказа (транзакционный запрос)
-- Обновляет статус отклика на "выбран" и устанавливает выбранного провайдера в заказе
-- Входные параметры: ID отклика
-- Возвращает: обновленную запись заказа
WITH selected_response AS (
    UPDATE order_responses
    SET 
        is_selected = true,
        updated_at = now()
    WHERE id = $1
    RETURNING *
)
UPDATE orders
SET 
    selected_provider_id = (SELECT provider_id FROM selected_response),
    status = 'accepted',
    updated_at = now()
WHERE id = (SELECT order_id FROM selected_response)
RETURNING *;

-- name: UnselectProviderForOrder :execrows
-- Отменяет выбор провайдера для заказа (транзакционный запрос)
-- Сбрасывает статус отклика и очищает поле выбранного провайдера в заказе
-- Входные параметры: ID отклика
-- Возвращает: количество обновленных строк (должно быть 1)
WITH unselect_response AS (
    UPDATE order_responses
    SET 
        is_selected = false,
        updated_at = now()
    WHERE id = $1
    RETURNING order_id
)
UPDATE orders
SET 
    selected_provider_id = NULL,
    status = 'pending',
    updated_at = now()
WHERE id = (SELECT order_id FROM unselect_response);

-- name: GetSelectedProviderForOrder :one
-- Получает данные выбранного провайдера для конкретного заказа
-- Входные параметры: ID заказа
-- Возвращает: запись пользователя-провайдера или ничего, если провайдер не выбран
SELECT u.* FROM users u
JOIN orders o ON u.id = o.selected_provider_id
WHERE o.id = $1
LIMIT 1;

-- name: GetSelectedResponseForOrder :one
-- Получает выбранный отклик для конкретного заказа
-- Входные параметры: ID заказа
-- Возвращает: запись выбранного отклика или ничего, если отклик не выбран
SELECT r.* FROM order_responses r
WHERE r.order_id = $1 AND r.is_selected = true
LIMIT 1;

-- name: GetOrdersWithSelectedProvider :many
-- Получает список заказов, для которых выбран провайдер, с детальной информацией
-- Входные параметры: лимит и смещение для пагинации
-- Возвращает: список заказов с данными о клиенте, категории и провайдере
SELECT o.*, u.username AS client_username, u.phone AS client_phone, u.whatsapp AS client_whatsapp,
       sc.name AS category_name, p.username AS provider_username
FROM orders o
JOIN users u ON o.client_id = u.id
JOIN service_categories sc ON o.category_id = sc.id
LEFT JOIN users p ON o.selected_provider_id = p.id
WHERE o.selected_provider_id IS NOT NULL
ORDER BY o.updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetOrdersWithResponsesByProvider :many
-- Получает все заказы, на которые провайдер оставил отклики (даже если не выбран)
-- Входные параметры: ID провайдера, лимит и смещение для пагинации
-- Возвращает: список заказов с базовой информацией о клиенте и категории
SELECT DISTINCT o.*, u.username AS client_username, sc.name AS category_name
FROM orders o
JOIN users u ON o.client_id = u.id
JOIN service_categories sc ON o.category_id = sc.id
JOIN order_responses r ON o.id = r.order_id
WHERE r.provider_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: HasProviderRespondedToOrder :one
-- Проверяет, откликнулся ли провайдер на заказ
-- Входные параметры: ID заказа, ID провайдера
-- Возвращает: булево значение (true/false)
SELECT EXISTS (
    SELECT 1 FROM order_responses
    WHERE order_id = $1 AND provider_id = $2
) AS has_responded;