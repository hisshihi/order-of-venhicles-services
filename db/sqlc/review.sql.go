// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: review.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const checkIfClientReviewedOrder = `-- name: CheckIfClientReviewedOrder :one
SELECT EXISTS(
        SELECT 1
        FROM "reviews"
        WHERE client_id = $1
            AND provider_id = $2
    ) as has_review
`

type CheckIfClientReviewedOrderParams struct {
	ClientID   int64 `json:"client_id"`
	ProviderID int64 `json:"provider_id"`
}

// Проверяет, оставил ли клиент отзыв по данному заказу
func (q *Queries) CheckIfClientReviewedOrder(ctx context.Context, arg CheckIfClientReviewedOrderParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkIfClientReviewedOrder, arg.ClientID, arg.ProviderID)
	var has_review bool
	err := row.Scan(&has_review)
	return has_review, err
}

const countReviews = `-- name: CountReviews :one
SELECT COUNT(*) FROM reviews
`

func (q *Queries) CountReviews(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countReviews)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createReview = `-- name: CreateReview :one
INSERT INTO "reviews" (
        client_id,
        provider_id,
        rating,
        comment
    )
VALUES ($1, $2, $3, $4)
RETURNING id, client_id, provider_id, rating, comment, created_at, updated_at
`

type CreateReviewParams struct {
	ClientID   int64  `json:"client_id"`
	ProviderID int64  `json:"provider_id"`
	Rating     int32  `json:"rating"`
	Comment    string `json:"comment"`
}

// Создает новый отзыв от клиента об услугодателе
func (q *Queries) CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error) {
	row := q.db.QueryRowContext(ctx, createReview,
		arg.ClientID,
		arg.ProviderID,
		arg.Rating,
		arg.Comment,
	)
	var i Review
	err := row.Scan(
		&i.ID,
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
DELETE FROM "reviews"
WHERE id = $1
    AND client_id = $2
`

type DeleteReviewParams struct {
	ID       int64 `json:"id"`
	ClientID int64 `json:"client_id"`
}

// Удаляет отзыв (только если пользователь является автором или администратором)
func (q *Queries) DeleteReview(ctx context.Context, arg DeleteReviewParams) error {
	_, err := q.db.ExecContext(ctx, deleteReview, arg.ID, arg.ClientID)
	return err
}

const getAverageRatingForProvider = `-- name: GetAverageRatingForProvider :one
SELECT COALESCE(AVG(rating), 0) as average_rating,
    COUNT(*) as total_reviews
FROM "reviews"
WHERE provider_id = $1
`

type GetAverageRatingForProviderRow struct {
	AverageRating sql.NullString `json:"average_rating"`
	TotalReviews  int64          `json:"total_reviews"`
}

// Получает среднюю оценку услугодателя
func (q *Queries) GetAverageRatingForProvider(ctx context.Context, providerID int64) (GetAverageRatingForProviderRow, error) {
	row := q.db.QueryRowContext(ctx, getAverageRatingForProvider, providerID)
	var i GetAverageRatingForProviderRow
	err := row.Scan(&i.AverageRating, &i.TotalReviews)
	return i, err
}

const getReviewByID = `-- name: GetReviewByID :one
SELECT r.id,
    r.client_id,
    r.provider_id,
    r.rating,
    r.comment,
    r.created_at,
    uc.username as client_name,
    up.username as provider_name
FROM "reviews" r
    JOIN "users" uc ON r.client_id = uc.id
    JOIN "users" up ON r.provider_id = up.id
WHERE r.id = $1
`

type GetReviewByIDRow struct {
	ID           int64     `json:"id"`
	ClientID     int64     `json:"client_id"`
	ProviderID   int64     `json:"provider_id"`
	Rating       int32     `json:"rating"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
	ClientName   string    `json:"client_name"`
	ProviderName string    `json:"provider_name"`
}

// Получает конкретный отзыв по ID
func (q *Queries) GetReviewByID(ctx context.Context, id int64) (GetReviewByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getReviewByID, id)
	var i GetReviewByIDRow
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.ProviderID,
		&i.Rating,
		&i.Comment,
		&i.CreatedAt,
		&i.ClientName,
		&i.ProviderName,
	)
	return i, err
}

const getReviewsByProviderID = `-- name: GetReviewsByProviderID :many
SELECT r.id,
    r.client_id,
    u.username as client_name,
    u.photo_url as client_photo,
    r.rating,
    r.comment,
    r.created_at
FROM "reviews" r
    JOIN "users" u ON r.client_id = u.id
WHERE r.provider_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3
`

type GetReviewsByProviderIDParams struct {
	ProviderID int64 `json:"provider_id"`
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
}

type GetReviewsByProviderIDRow struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"client_id"`
	ClientName  string    `json:"client_name"`
	ClientPhoto []byte    `json:"client_photo"`
	Rating      int32     `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

// Получает все отзывы об услугодателе
func (q *Queries) GetReviewsByProviderID(ctx context.Context, arg GetReviewsByProviderIDParams) ([]GetReviewsByProviderIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getReviewsByProviderID, arg.ProviderID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetReviewsByProviderIDRow{}
	for rows.Next() {
		var i GetReviewsByProviderIDRow
		if err := rows.Scan(
			&i.ID,
			&i.ClientID,
			&i.ClientName,
			&i.ClientPhoto,
			&i.Rating,
			&i.Comment,
			&i.CreatedAt,
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

const listReview = `-- name: ListReview :many
SELECT id, client_id, provider_id, rating, comment, created_at, updated_at FROM "reviews" r
ORDER BY r.created_at
LIMIT $1 OFFSET $2
`

type ListReviewParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListReview(ctx context.Context, arg ListReviewParams) ([]Review, error) {
	rows, err := q.db.QueryContext(ctx, listReview, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Review{}
	for rows.Next() {
		var i Review
		if err := rows.Scan(
			&i.ID,
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
