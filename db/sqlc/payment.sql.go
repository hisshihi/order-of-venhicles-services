// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: payment.sql

package sqlc

import (
	"context"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (user_id, amount, payment_method, status)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, amount, payment_method, status, created_at, updated_at
`

type CreatePaymentParams struct {
	UserID        int64             `json:"user_id"`
	Amount        string            `json:"amount"`
	PaymentMethod string            `json:"payment_method"`
	Status        NullStatusPayment `json:"status"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.UserID,
		arg.Amount,
		arg.PaymentMethod,
		arg.Status,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.PaymentMethod,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1
`

func (q *Queries) DeletePayment(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePayment, id)
	return err
}

const getPaymentByID = `-- name: GetPaymentByID :one
SELECT id, user_id, amount, payment_method, status, created_at, updated_at FROM payments
WHERE id = $1
`

func (q *Queries) GetPaymentByID(ctx context.Context, id int64) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPaymentByID, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.PaymentMethod,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPayments = `-- name: ListPayments :many
SELECT id, user_id, amount, payment_method, status, created_at, updated_at FROM payments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListPaymentsParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListPayments(ctx context.Context, arg ListPaymentsParams) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, listPayments, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Payment{}
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Amount,
			&i.PaymentMethod,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updatePayment = `-- name: UpdatePayment :one
UPDATE payments
SET user_id = $2, amount = $3, payment_method = $4, status = $5
WHERE id = $1
RETURNING id, user_id, amount, payment_method, status, created_at, updated_at
`

type UpdatePaymentParams struct {
	ID            int64             `json:"id"`
	UserID        int64             `json:"user_id"`
	Amount        string            `json:"amount"`
	PaymentMethod string            `json:"payment_method"`
	Status        NullStatusPayment `json:"status"`
}

func (q *Queries) UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, updatePayment,
		arg.ID,
		arg.UserID,
		arg.Amount,
		arg.PaymentMethod,
		arg.Status,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.PaymentMethod,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
