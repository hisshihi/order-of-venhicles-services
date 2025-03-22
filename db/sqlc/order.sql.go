// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: order.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const acceptOrderByProviderID = `-- name: AcceptOrderByProviderID :one
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
RETURNING id, client_id, category_id, service_id, status, created_at, updated_at, provider_accepted, provider_message, client_message, order_date, selected_provider_id
`

type AcceptOrderByProviderIDParams struct {
	Column1 sql.NullInt64  `json:"column_1"`
	Column2 sql.NullInt64  `json:"column_2"`
	Column3 sql.NullInt64  `json:"column_3"`
	Column4 sql.NullString `json:"column_4"`
}

// Провайдер принимает заказ и указывает свою услугу
func (q *Queries) AcceptOrderByProviderID(ctx context.Context, arg AcceptOrderByProviderIDParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, acceptOrderByProviderID,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Column4,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.CategoryID,
		&i.ServiceID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ProviderAccepted,
		&i.ProviderMessage,
		&i.ClientMessage,
		&i.OrderDate,
		&i.SelectedProviderID,
	)
	return i, err
}

const createOrder = `-- name: CreateOrder :one
INSERT INTO "orders" (
        client_id,
        category_id,
        status,
        client_message,
        order_date
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING id, client_id, category_id, service_id, status, created_at, updated_at, provider_accepted, provider_message, client_message, order_date, selected_provider_id
`

type CreateOrderParams struct {
	ClientID      int64            `json:"client_id"`
	CategoryID    int64            `json:"category_id"`
	Status        NullStatusOrders `json:"status"`
	ClientMessage sql.NullString   `json:"client_message"`
	OrderDate     sql.NullTime     `json:"order_date"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.ClientID,
		arg.CategoryID,
		arg.Status,
		arg.ClientMessage,
		arg.OrderDate,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.CategoryID,
		&i.ServiceID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ProviderAccepted,
		&i.ProviderMessage,
		&i.ClientMessage,
		&i.OrderDate,
		&i.SelectedProviderID,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1
`

func (q *Queries) DeleteOrder(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteOrder, id)
	return err
}

const getOrderByID = `-- name: GetOrderByID :one
SELECT o.id, o.client_id, o.category_id, o.service_id, o.status, o.created_at, o.updated_at, o.provider_accepted, o.provider_message, o.client_message, o.order_date, o.selected_provider_id,
    sc.name as category_name,
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
    JOIN "users" u ON o.client_id = u.id
    LEFT JOIN "services" s ON s.id = o.service_id
    LEFT JOIN "users" p ON s.provider_id = p.id
WHERE o.id = $1
`

type GetOrderByIDRow struct {
	ID                 int64            `json:"id"`
	ClientID           int64            `json:"client_id"`
	CategoryID         int64            `json:"category_id"`
	ServiceID          sql.NullInt64    `json:"service_id"`
	Status             NullStatusOrders `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	ProviderAccepted   sql.NullBool     `json:"provider_accepted"`
	ProviderMessage    sql.NullString   `json:"provider_message"`
	ClientMessage      sql.NullString   `json:"client_message"`
	OrderDate          sql.NullTime     `json:"order_date"`
	SelectedProviderID sql.NullInt64    `json:"selected_provider_id"`
	CategoryName       string           `json:"category_name"`
	ClientName         string           `json:"client_name"`
	ClientPhone        string           `json:"client_phone"`
	ClientWhatsapp     string           `json:"client_whatsapp"`
	ClientCity         sql.NullString   `json:"client_city"`
	ClientDistrict     sql.NullString   `json:"client_district"`
	ProviderName       sql.NullString   `json:"provider_name"`
	ProviderPhone      sql.NullString   `json:"provider_phone"`
	ProviderWhatsapp   sql.NullString   `json:"provider_whatsapp"`
	ServiceTitle       sql.NullString   `json:"service_title"`
}

func (q *Queries) GetOrderByID(ctx context.Context, id int64) (GetOrderByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getOrderByID, id)
	var i GetOrderByIDRow
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.CategoryID,
		&i.ServiceID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ProviderAccepted,
		&i.ProviderMessage,
		&i.ClientMessage,
		&i.OrderDate,
		&i.SelectedProviderID,
		&i.CategoryName,
		&i.ClientName,
		&i.ClientPhone,
		&i.ClientWhatsapp,
		&i.ClientCity,
		&i.ClientDistrict,
		&i.ProviderName,
		&i.ProviderPhone,
		&i.ProviderWhatsapp,
		&i.ServiceTitle,
	)
	return i, err
}

const getOrderStatistics = `-- name: GetOrderStatistics :one
WITH provider_orders AS (
    SELECT o.id, o.client_id, o.category_id, o.service_id, o.status, o.created_at, o.updated_at, o.provider_accepted, o.provider_message, o.client_message, o.order_date, o.selected_provider_id
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
FROM provider_orders
`

type GetOrderStatisticsRow struct {
	PendingCount   int64 `json:"pending_count"`
	AcceptedCount  int64 `json:"accepted_count"`
	CompletedCount int64 `json:"completed_count"`
	CancelledCount int64 `json:"cancelled_count"`
	TotalCount     int64 `json:"total_count"`
}

// Получает статистику заказов для услугодателя
func (q *Queries) GetOrderStatistics(ctx context.Context, dollar_1 sql.NullInt64) (GetOrderStatisticsRow, error) {
	row := q.db.QueryRowContext(ctx, getOrderStatistics, dollar_1)
	var i GetOrderStatisticsRow
	err := row.Scan(
		&i.PendingCount,
		&i.AcceptedCount,
		&i.CompletedCount,
		&i.CancelledCount,
		&i.TotalCount,
	)
	return i, err
}

const getOrdersByCategory = `-- name: GetOrdersByCategory :many
SELECT o.id, o.client_id, o.category_id, o.service_id, o.status, o.created_at, o.updated_at, o.provider_accepted, o.provider_message, o.client_message, o.order_date, o.selected_provider_id,
    sc.name as category_name,
    u.username as client_name
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
    JOIN "users" u ON o.client_id = u.id
WHERE o.category_id = $1
    AND o.status = 'pending'
    AND o.provider_accepted = false
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3
`

type GetOrdersByCategoryParams struct {
	CategoryID int64 `json:"category_id"`
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
}

type GetOrdersByCategoryRow struct {
	ID                 int64            `json:"id"`
	ClientID           int64            `json:"client_id"`
	CategoryID         int64            `json:"category_id"`
	ServiceID          sql.NullInt64    `json:"service_id"`
	Status             NullStatusOrders `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	ProviderAccepted   sql.NullBool     `json:"provider_accepted"`
	ProviderMessage    sql.NullString   `json:"provider_message"`
	ClientMessage      sql.NullString   `json:"client_message"`
	OrderDate          sql.NullTime     `json:"order_date"`
	SelectedProviderID sql.NullInt64    `json:"selected_provider_id"`
	CategoryName       string           `json:"category_name"`
	ClientName         string           `json:"client_name"`
}

// Получает список заказов по категории
func (q *Queries) GetOrdersByCategory(ctx context.Context, arg GetOrdersByCategoryParams) ([]GetOrdersByCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrdersByCategory, arg.CategoryID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOrdersByCategoryRow{}
	for rows.Next() {
		var i GetOrdersByCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.ClientID,
			&i.CategoryID,
			&i.ServiceID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ProviderAccepted,
			&i.ProviderMessage,
			&i.ClientMessage,
			&i.OrderDate,
			&i.SelectedProviderID,
			&i.CategoryName,
			&i.ClientName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAvailableOrdersForProvider = `-- name: ListAvailableOrdersForProvider :many
SELECT o.id, o.client_id, o.category_id, o.service_id, o.status, o.created_at, o.updated_at, o.provider_accepted, o.provider_message, o.client_message, o.order_date, o.selected_provider_id,
    sc.name as category_name,
    u.username as client_name,
    u.city as client_city,
    u.district as client_district
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
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
LIMIT $1 OFFSET $2
`

type ListAvailableOrdersForProviderParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type ListAvailableOrdersForProviderRow struct {
	ID                 int64            `json:"id"`
	ClientID           int64            `json:"client_id"`
	CategoryID         int64            `json:"category_id"`
	ServiceID          sql.NullInt64    `json:"service_id"`
	Status             NullStatusOrders `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	ProviderAccepted   sql.NullBool     `json:"provider_accepted"`
	ProviderMessage    sql.NullString   `json:"provider_message"`
	ClientMessage      sql.NullString   `json:"client_message"`
	OrderDate          sql.NullTime     `json:"order_date"`
	SelectedProviderID sql.NullInt64    `json:"selected_provider_id"`
	CategoryName       string           `json:"category_name"`
	ClientName         string           `json:"client_name"`
	ClientCity         sql.NullString   `json:"client_city"`
	ClientDistrict     sql.NullString   `json:"client_district"`
}

// Получает список доступных заказов для провайдера услуг
func (q *Queries) ListAvailableOrdersForProvider(ctx context.Context, arg ListAvailableOrdersForProviderParams) ([]ListAvailableOrdersForProviderRow, error) {
	rows, err := q.db.QueryContext(ctx, listAvailableOrdersForProvider, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAvailableOrdersForProviderRow{}
	for rows.Next() {
		var i ListAvailableOrdersForProviderRow
		if err := rows.Scan(
			&i.ID,
			&i.ClientID,
			&i.CategoryID,
			&i.ServiceID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ProviderAccepted,
			&i.ProviderMessage,
			&i.ClientMessage,
			&i.OrderDate,
			&i.SelectedProviderID,
			&i.CategoryName,
			&i.ClientName,
			&i.ClientCity,
			&i.ClientDistrict,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCountAvailableOrdersForProvider = `-- name: ListCountAvailableOrdersForProvider :one
SELECT COUNT(*) FROM "orders"
WHERE status = 'pending'
    AND provider_accepted = false
    AND category_id IN (
        SELECT DISTINCT category_id
        FROM "services"
        WHERE provider_id = $1
    )
`

func (q *Queries) ListCountAvailableOrdersForProvider(ctx context.Context, providerID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, listCountAvailableOrdersForProvider, providerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listCountOrdersByClientID = `-- name: ListCountOrdersByClientID :one
SELECT COUNT(*) FROM "orders"
WHERE client_id = $1
`

func (q *Queries) ListCountOrdersByClientID(ctx context.Context, clientID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, listCountOrdersByClientID, clientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listOrdersByClientID = `-- name: ListOrdersByClientID :many
SELECT o.id, o.client_id, o.category_id, o.service_id, o.status, o.created_at, o.updated_at, o.provider_accepted, o.provider_message, o.client_message, o.order_date, o.selected_provider_id,
    sc.name as category_name,
    s.title as service_title,
    p.username as provider_name
FROM "orders" o
    JOIN "service_categories" sc ON o.category_id = sc.id
    LEFT JOIN "services" s ON o.service_id = s.id
    LEFT JOIN "users" p ON s.provider_id = p.id
WHERE o.client_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3
`

type ListOrdersByClientIDParams struct {
	ClientID int64 `json:"client_id"`
	Limit    int64 `json:"limit"`
	Offset   int64 `json:"offset"`
}

type ListOrdersByClientIDRow struct {
	ID                 int64            `json:"id"`
	ClientID           int64            `json:"client_id"`
	CategoryID         int64            `json:"category_id"`
	ServiceID          sql.NullInt64    `json:"service_id"`
	Status             NullStatusOrders `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	ProviderAccepted   sql.NullBool     `json:"provider_accepted"`
	ProviderMessage    sql.NullString   `json:"provider_message"`
	ClientMessage      sql.NullString   `json:"client_message"`
	OrderDate          sql.NullTime     `json:"order_date"`
	SelectedProviderID sql.NullInt64    `json:"selected_provider_id"`
	CategoryName       string           `json:"category_name"`
	ServiceTitle       sql.NullString   `json:"service_title"`
	ProviderName       sql.NullString   `json:"provider_name"`
}

func (q *Queries) ListOrdersByClientID(ctx context.Context, arg ListOrdersByClientIDParams) ([]ListOrdersByClientIDRow, error) {
	rows, err := q.db.QueryContext(ctx, listOrdersByClientID, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListOrdersByClientIDRow{}
	for rows.Next() {
		var i ListOrdersByClientIDRow
		if err := rows.Scan(
			&i.ID,
			&i.ClientID,
			&i.CategoryID,
			&i.ServiceID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ProviderAccepted,
			&i.ProviderMessage,
			&i.ClientMessage,
			&i.OrderDate,
			&i.SelectedProviderID,
			&i.CategoryName,
			&i.ServiceTitle,
			&i.ProviderName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrder = `-- name: UpdateOrder :one
UPDATE orders
SET category_id = $2,
    client_message = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, client_id, category_id, service_id, status, created_at, updated_at, provider_accepted, provider_message, client_message, order_date, selected_provider_id
`

type UpdateOrderParams struct {
	ID            int64          `json:"id"`
	CategoryID    int64          `json:"category_id"`
	ClientMessage sql.NullString `json:"client_message"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrder, arg.ID, arg.CategoryID, arg.ClientMessage)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.CategoryID,
		&i.ServiceID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ProviderAccepted,
		&i.ProviderMessage,
		&i.ClientMessage,
		&i.OrderDate,
		&i.SelectedProviderID,
	)
	return i, err
}

const updateOrderStatus = `-- name: UpdateOrderStatus :one
UPDATE "orders"
SET status = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, client_id, category_id, service_id, status, created_at, updated_at, provider_accepted, provider_message, client_message, order_date, selected_provider_id
`

type UpdateOrderStatusParams struct {
	ID     int64            `json:"id"`
	Status NullStatusOrders `json:"status"`
}

func (q *Queries) UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrderStatus, arg.ID, arg.Status)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.CategoryID,
		&i.ServiceID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ProviderAccepted,
		&i.ProviderMessage,
		&i.ClientMessage,
		&i.OrderDate,
		&i.SelectedProviderID,
	)
	return i, err
}
