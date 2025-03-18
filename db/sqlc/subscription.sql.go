// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: subscription.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const checkAndUpdateExpiredSubscriptions = `-- name: CheckAndUpdateExpiredSubscriptions :many
UPDATE subscriptions
SET status = 'expired',
    updated_at = NOW()
WHERE status = 'active'
    AND end_date <= NOW()
RETURNING id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
`

func (q *Queries) CheckAndUpdateExpiredSubscriptions(ctx context.Context) ([]Subscription, error) {
	rows, err := q.db.QueryContext(ctx, checkAndUpdateExpiredSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subscription{}
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.StartDate,
			&i.EndDate,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SubscriptionType,
			&i.Price,
			&i.PromoCodeID,
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

const createSubscription = `-- name: CreateSubscription :one
INSERT INTO subscriptions (provider_id, start_date, subscription_type, promo_code_id, price, end_date, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
`

type CreateSubscriptionParams struct {
	ProviderID       int64                  `json:"provider_id"`
	StartDate        time.Time              `json:"start_date"`
	SubscriptionType sql.NullString         `json:"subscription_type"`
	PromoCodeID      sql.NullInt64          `json:"promo_code_id"`
	Price            sql.NullString         `json:"price"`
	EndDate          time.Time              `json:"end_date"`
	Status           NullStatusSubscription `json:"status"`
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, createSubscription,
		arg.ProviderID,
		arg.StartDate,
		arg.SubscriptionType,
		arg.PromoCodeID,
		arg.Price,
		arg.EndDate,
		arg.Status,
	)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.StartDate,
		&i.EndDate,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SubscriptionType,
		&i.Price,
		&i.PromoCodeID,
	)
	return i, err
}

const deleteSubscription = `-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE id = $1
`

func (q *Queries) DeleteSubscription(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSubscription, id)
	return err
}

const getActiveSubscriptionForProvider = `-- name: GetActiveSubscriptionForProvider :one
SELECT id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
FROM subscriptions
WHERE provider_id = $1
    AND status = 'active'
    AND end_date > NOW()
ORDER BY end_date DESC
LIMIT 1
`

func (q *Queries) GetActiveSubscriptionForProvider(ctx context.Context, providerID int64) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, getActiveSubscriptionForProvider, providerID)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.StartDate,
		&i.EndDate,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SubscriptionType,
		&i.Price,
		&i.PromoCodeID,
	)
	return i, err
}

const getSubscriptionByID = `-- name: GetSubscriptionByID :one
SELECT id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
FROM subscriptions
WHERE id = $1
`

func (q *Queries) GetSubscriptionByID(ctx context.Context, id int64) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, getSubscriptionByID, id)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.StartDate,
		&i.EndDate,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SubscriptionType,
		&i.Price,
		&i.PromoCodeID,
	)
	return i, err
}

const getSubscriptionByProviderID = `-- name: GetSubscriptionByProviderID :one
SELECT id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
FROM subscriptions
WHERE provider_id = $1
`

func (q *Queries) GetSubscriptionByProviderID(ctx context.Context, providerID int64) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, getSubscriptionByProviderID, providerID)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.StartDate,
		&i.EndDate,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SubscriptionType,
		&i.Price,
		&i.PromoCodeID,
	)
	return i, err
}

const listSubscriptions = `-- name: ListSubscriptions :many
SELECT id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
FROM subscriptions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListSubscriptionsParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListSubscriptions(ctx context.Context, arg ListSubscriptionsParams) ([]Subscription, error) {
	rows, err := q.db.QueryContext(ctx, listSubscriptions, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subscription{}
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.StartDate,
			&i.EndDate,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SubscriptionType,
			&i.Price,
			&i.PromoCodeID,
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

const listSubscriptionsByProviderID = `-- name: ListSubscriptionsByProviderID :many
SELECT id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
FROM subscriptions
WHERE provider_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type ListSubscriptionsByProviderIDParams struct {
	ProviderID int64 `json:"provider_id"`
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
}

func (q *Queries) ListSubscriptionsByProviderID(ctx context.Context, arg ListSubscriptionsByProviderIDParams) ([]Subscription, error) {
	rows, err := q.db.QueryContext(ctx, listSubscriptionsByProviderID, arg.ProviderID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subscription{}
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.StartDate,
			&i.EndDate,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SubscriptionType,
			&i.Price,
			&i.PromoCodeID,
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

const updateSubscription = `-- name: UpdateSubscription :one
UPDATE subscriptions
SET provider_id = $2,
    start_date = $3,
    end_date = $4,
    status = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING id, provider_id, start_date, end_date, status, created_at, updated_at, subscription_type, price, promo_code_id
`

type UpdateSubscriptionParams struct {
	ID         int64                  `json:"id"`
	ProviderID int64                  `json:"provider_id"`
	StartDate  time.Time              `json:"start_date"`
	EndDate    time.Time              `json:"end_date"`
	Status     NullStatusSubscription `json:"status"`
}

func (q *Queries) UpdateSubscription(ctx context.Context, arg UpdateSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, updateSubscription,
		arg.ID,
		arg.ProviderID,
		arg.StartDate,
		arg.EndDate,
		arg.Status,
	)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.StartDate,
		&i.EndDate,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SubscriptionType,
		&i.Price,
		&i.PromoCodeID,
	)
	return i, err
}
