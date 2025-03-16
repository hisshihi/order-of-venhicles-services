// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: service.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createService = `-- name: CreateService :one
INSERT INTO services (provider_id, category_id, title, description, price)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, provider_id, category_id, title, description, price, created_at, updated_at
`

type CreateServiceParams struct {
	ProviderID  int64  `json:"provider_id"`
	CategoryID  int64  `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func (q *Queries) CreateService(ctx context.Context, arg CreateServiceParams) (Service, error) {
	row := q.db.QueryRowContext(ctx, createService,
		arg.ProviderID,
		arg.CategoryID,
		arg.Title,
		arg.Description,
		arg.Price,
	)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.CategoryID,
		&i.Title,
		&i.Description,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteService = `-- name: DeleteService :exec
DELETE FROM services
WHERE id = $1
`

func (q *Queries) DeleteService(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteService, id)
	return err
}

const getServiceByID = `-- name: GetServiceByID :one
SELECT id, provider_id, category_id, title, description, price, created_at, updated_at FROM services
WHERE id = $1
`

func (q *Queries) GetServiceByID(ctx context.Context, id int64) (Service, error) {
	row := q.db.QueryRowContext(ctx, getServiceByID, id)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.CategoryID,
		&i.Title,
		&i.Description,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getServiceByProviderID = `-- name: GetServiceByProviderID :one
SELECT id, provider_id, category_id, title, description, price, created_at, updated_at FROM services
WHERE provider_id = $1
`

func (q *Queries) GetServiceByProviderID(ctx context.Context, providerID int64) (Service, error) {
	row := q.db.QueryRowContext(ctx, getServiceByProviderID, providerID)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.CategoryID,
		&i.Title,
		&i.Description,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listServices = `-- name: ListServices :many
SELECT id, provider_id, category_id, title, description, price, created_at, updated_at FROM services
ORDER BY title DESC
LIMIT $1 OFFSET $2
`

type ListServicesParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListServices(ctx context.Context, arg ListServicesParams) ([]Service, error) {
	rows, err := q.db.QueryContext(ctx, listServices, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Service{}
	for rows.Next() {
		var i Service
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.CategoryID,
			&i.Title,
			&i.Description,
			&i.Price,
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

const listServicesByCategory = `-- name: ListServicesByCategory :many
SELECT id, provider_id, category_id, title, description, price, created_at, updated_at FROM services
WHERE category_id = $1
ORDER BY title
LIMIT $2 OFFSET $3
`

type ListServicesByCategoryParams struct {
	CategoryID int64 `json:"category_id"`
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
}

func (q *Queries) ListServicesByCategory(ctx context.Context, arg ListServicesByCategoryParams) ([]Service, error) {
	rows, err := q.db.QueryContext(ctx, listServicesByCategory, arg.CategoryID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Service{}
	for rows.Next() {
		var i Service
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.CategoryID,
			&i.Title,
			&i.Description,
			&i.Price,
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

const listServicesByProviderID = `-- name: ListServicesByProviderID :many
SELECT id, provider_id, category_id, title, description, price, created_at, updated_at FROM services
WHERE provider_id = $1
ORDER BY title
LIMIT $2 OFFSET $3
`

type ListServicesByProviderIDParams struct {
	ProviderID int64 `json:"provider_id"`
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
}

func (q *Queries) ListServicesByProviderID(ctx context.Context, arg ListServicesByProviderIDParams) ([]Service, error) {
	rows, err := q.db.QueryContext(ctx, listServicesByProviderID, arg.ProviderID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Service{}
	for rows.Next() {
		var i Service
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.CategoryID,
			&i.Title,
			&i.Description,
			&i.Price,
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

const listServicesByTitle = `-- name: ListServicesByTitle :many
SELECT id, provider_id, category_id, title, description, price, created_at, updated_at FROM services
WHERE title ILIKE '%' || $1 || '%'
ORDER BY title
LIMIT $2 OFFSET $3
`

type ListServicesByTitleParams struct {
	Column1 sql.NullString `json:"column_1"`
	Limit   int64          `json:"limit"`
	Offset  int64          `json:"offset"`
}

func (q *Queries) ListServicesByTitle(ctx context.Context, arg ListServicesByTitleParams) ([]Service, error) {
	rows, err := q.db.QueryContext(ctx, listServicesByTitle, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Service{}
	for rows.Next() {
		var i Service
		if err := rows.Scan(
			&i.ID,
			&i.ProviderID,
			&i.CategoryID,
			&i.Title,
			&i.Description,
			&i.Price,
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

const updateService = `-- name: UpdateService :one
UPDATE services
SET provider_id = $2, category_id = $3, title = $4, description = $5, price = $6, updated_at = NOW()
WHERE id = $1
RETURNING id, provider_id, category_id, title, description, price, created_at, updated_at
`

type UpdateServiceParams struct {
	ID          int64  `json:"id"`
	ProviderID  int64  `json:"provider_id"`
	CategoryID  int64  `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func (q *Queries) UpdateService(ctx context.Context, arg UpdateServiceParams) (Service, error) {
	row := q.db.QueryRowContext(ctx, updateService,
		arg.ID,
		arg.ProviderID,
		arg.CategoryID,
		arg.Title,
		arg.Description,
		arg.Price,
	)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.CategoryID,
		&i.Title,
		&i.Description,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
