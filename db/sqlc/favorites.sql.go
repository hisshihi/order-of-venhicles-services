// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: favorites.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const addProviderToFavorites = `-- name: AddProviderToFavorites :one
INSERT INTO "favorites" (client_id, provider_id)
VALUES ($1, $2) ON CONFLICT (client_id, provider_id) DO NOTHING
RETURNING id, client_id, provider_id, created_at
`

type AddProviderToFavoritesParams struct {
	ClientID   int64 `json:"client_id"`
	ProviderID int64 `json:"provider_id"`
}

// Добавляет услугодателя в избранное клиента
func (q *Queries) AddProviderToFavorites(ctx context.Context, arg AddProviderToFavoritesParams) (Favorite, error) {
	row := q.db.QueryRowContext(ctx, addProviderToFavorites, arg.ClientID, arg.ProviderID)
	var i Favorite
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.ProviderID,
		&i.CreatedAt,
	)
	return i, err
}

const checkIfProviderIsFavorite = `-- name: CheckIfProviderIsFavorite :one
SELECT EXISTS(
        SELECT 1
        FROM "favorites"
        WHERE client_id = $1
            AND provider_id = $2
    ) as is_favorite
`

type CheckIfProviderIsFavoriteParams struct {
	ClientID   int64 `json:"client_id"`
	ProviderID int64 `json:"provider_id"`
}

// Проверяет, добавлен ли услугодатель в избранное клиента
func (q *Queries) CheckIfProviderIsFavorite(ctx context.Context, arg CheckIfProviderIsFavoriteParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkIfProviderIsFavorite, arg.ClientID, arg.ProviderID)
	var is_favorite bool
	err := row.Scan(&is_favorite)
	return is_favorite, err
}

const listFavoriteProviders = `-- name: ListFavoriteProviders :many
SELECT f.id,
    f.provider_id,
    u.username as provider_name,
    u.photo_url as provider_photo,
    u.phone as provider_phone,
    u.whatsapp as provider_whatsapp,
    u.description as provider_description,
    (
        SELECT COUNT(*)
        FROM "reviews" r
        WHERE r.provider_id = f.provider_id
    ) as reviews_count,
    (
        SELECT COALESCE(AVG(r.rating), 0)
        FROM "reviews" r
        WHERE r.provider_id = f.provider_id
    ) as average_rating,
    f.created_at
FROM "favorites" f
    JOIN "users" u ON f.provider_id = u.id
WHERE f.client_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3
`

type ListFavoriteProvidersParams struct {
	ClientID int64 `json:"client_id"`
	Limit    int64 `json:"limit"`
	Offset   int64 `json:"offset"`
}

type ListFavoriteProvidersRow struct {
	ID                  int64          `json:"id"`
	ProviderID          int64          `json:"provider_id"`
	ProviderName        string         `json:"provider_name"`
	ProviderPhoto       []byte         `json:"provider_photo"`
	ProviderPhone       string         `json:"provider_phone"`
	ProviderWhatsapp    string         `json:"provider_whatsapp"`
	ProviderDescription sql.NullString `json:"provider_description"`
	ReviewsCount        int64          `json:"reviews_count"`
	AverageRating       sql.NullString `json:"average_rating"`
	CreatedAt           time.Time      `json:"created_at"`
}

// Получает список избранных услугодателей клиента
func (q *Queries) ListFavoriteProviders(ctx context.Context, arg ListFavoriteProvidersParams) ([]ListFavoriteProvidersRow, error) {
	rows, err := q.db.QueryContext(ctx, listFavoriteProviders, arg.ClientID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListFavoriteProvidersRow{}
	for rows.Next() {
		var i ListFavoriteProvidersRow
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.ProviderName,
			&i.ProviderPhoto,
			&i.ProviderPhone,
			&i.ProviderWhatsapp,
			&i.ProviderDescription,
			&i.ReviewsCount,
			&i.AverageRating,
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

const removeProviderFromFavorites = `-- name: RemoveProviderFromFavorites :exec
DELETE FROM "favorites"
WHERE client_id = $1
    AND provider_id = $2
`

type RemoveProviderFromFavoritesParams struct {
	ClientID   int64 `json:"client_id"`
	ProviderID int64 `json:"provider_id"`
}

// Удаляет услугодателя из избранного клиента
func (q *Queries) RemoveProviderFromFavorites(ctx context.Context, arg RemoveProviderFromFavoritesParams) error {
	_, err := q.db.ExecContext(ctx, removeProviderFromFavorites, arg.ClientID, arg.ProviderID)
	return err
}
