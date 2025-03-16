// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: review.sql

package sqlc

import (
	"context"
)

const createReview = `-- name: CreateReview :one
INSERT INTO reviews (order_id, client_id, provider_id, rating, comment)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, order_id, client_id, provider_id, rating, comment, created_at, updated_at
`

type CreateReviewParams struct {
	OrderID    int64  `json:"order_id"`
	ClientID   int64  `json:"client_id"`
	ProviderID int64  `json:"provider_id"`
	Rating     int32  `json:"rating"`
	Comment    string `json:"comment"`
}

func (q *Queries) CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error) {
	row := q.db.QueryRowContext(ctx, createReview,
		arg.OrderID,
		arg.ClientID,
		arg.ProviderID,
		arg.Rating,
		arg.Comment,
	)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ClientID,
		&i.ProviderID,
		&i.Rating,
		&i.Comment,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteReview = `-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1
`

func (q *Queries) DeleteReview(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteReview, id)
	return err
}

const getProviderOverallRating = `-- name: GetProviderOverallRating :one
SELECT provider_id, AVG(rating) AS overall_rating, COUNT(*) AS review_count
FROM reviews
WHERE provider_id = $1
GROUP BY provider_id
`

type GetProviderOverallRatingRow struct {
	ProviderID    int64  `json:"provider_id"`
	OverallRating string `json:"overall_rating"`
	ReviewCount   int64  `json:"review_count"`
}

func (q *Queries) GetProviderOverallRating(ctx context.Context, providerID int64) (GetProviderOverallRatingRow, error) {
	row := q.db.QueryRowContext(ctx, getProviderOverallRating, providerID)
	var i GetProviderOverallRatingRow
	err := row.Scan(&i.ProviderID, &i.OverallRating, &i.ReviewCount)
	return i, err
}

const getReviewByID = `-- name: GetReviewByID :one
SELECT id, order_id, client_id, provider_id, rating, comment, created_at, updated_at FROM reviews
WHERE id = $1
`

func (q *Queries) GetReviewByID(ctx context.Context, id int64) (Review, error) {
	row := q.db.QueryRowContext(ctx, getReviewByID, id)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ClientID,
		&i.ProviderID,
		&i.Rating,
		&i.Comment,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listReviewsByProviderID = `-- name: ListReviewsByProviderID :many
SELECT id, order_id, client_id, provider_id, rating, comment, created_at, updated_at FROM reviews
WHERE provider_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type ListReviewsByProviderIDParams struct {
	ProviderID int64 `json:"provider_id"`
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
}

func (q *Queries) ListReviewsByProviderID(ctx context.Context, arg ListReviewsByProviderIDParams) ([]Review, error) {
	rows, err := q.db.QueryContext(ctx, listReviewsByProviderID, arg.ProviderID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Review{}
	for rows.Next() {
		var i Review
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ClientID,
			&i.ProviderID,
			&i.Rating,
			&i.Comment,
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

const updateReview = `-- name: UpdateReview :one
UPDATE reviews
SET rating = $2, comment = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, order_id, client_id, provider_id, rating, comment, created_at, updated_at
`

type UpdateReviewParams struct {
	ID      int64  `json:"id"`
	Rating  int32  `json:"rating"`
	Comment string `json:"comment"`
}

func (q *Queries) UpdateReview(ctx context.Context, arg UpdateReviewParams) (Review, error) {
	row := q.db.QueryRowContext(ctx, updateReview, arg.ID, arg.Rating, arg.Comment)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ClientID,
		&i.ProviderID,
		&i.Rating,
		&i.Comment,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
