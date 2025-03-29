// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: city.sql

package sqlc

import (
	"context"
)

const createCity = `-- name: CreateCity :one
INSERT INTO city (name)
VALUES ($1)
RETURNING id, name
`

func (q *Queries) CreateCity(ctx context.Context, name string) (City, error) {
	row := q.db.QueryRowContext(ctx, createCity, name)
	var i City
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteCity = `-- name: DeleteCity :exec
DELETE FROM city
WHERE id = $1
`

func (q *Queries) DeleteCity(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCity, id)
	return err
}

const getCityByID = `-- name: GetCityByID :one
SELECT id, name FROM city
WHERE id = $1
`

func (q *Queries) GetCityByID(ctx context.Context, id int64) (City, error) {
	row := q.db.QueryRowContext(ctx, getCityByID, id)
	var i City
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listCity = `-- name: ListCity :many
SELECT id, name FROM city
ORDER BY name
`

func (q *Queries) ListCity(ctx context.Context) ([]City, error) {
	rows, err := q.db.QueryContext(ctx, listCity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []City{}
	for rows.Next() {
		var i City
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const updateCity = `-- name: UpdateCity :one
UPDATE city
SET name = $2
WHERE id = $1
RETURNING id, name
`

type UpdateCityParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateCity(ctx context.Context, arg UpdateCityParams) (City, error) {
	row := q.db.QueryRowContext(ctx, updateCity, arg.ID, arg.Name)
	var i City
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
