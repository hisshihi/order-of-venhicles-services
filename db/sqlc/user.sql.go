// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package sqlc

import (
	"context"
	"database/sql"
)

const blockedUser = `-- name: BlockedUser :one
UPDATE users
SET is_blocked = true
WHERE id = $1
RETURNING is_blocked
`

func (q *Queries) BlockedUser(ctx context.Context, id int64) (sql.NullBool, error) {
	row := q.db.QueryRowContext(ctx, blockedUser, id)
	var is_blocked sql.NullBool
	err := row.Scan(&is_blocked)
	return is_blocked, err
}

const changePassword = `-- name: ChangePassword :exec
UPDATE users
SET password_hash = $2,
    password_change_at = NOW()
WHERE id = $1
`

type ChangePasswordParams struct {
	ID           int64  `json:"id"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) ChangePassword(ctx context.Context, arg ChangePasswordParams) error {
	_, err := q.db.ExecContext(ctx, changePassword, arg.ID, arg.PasswordHash)
	return err
}

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*)
FROM users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
        username,
        email,
        password_hash,
        role,
        country,
        city,
        district,
        phone,
        whatsapp,
        photo_url
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10
    )
RETURNING id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
`

type CreateUserParams struct {
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"password_hash"`
	Role         NullRole       `json:"role"`
	Country      sql.NullString `json:"country"`
	City         sql.NullString `json:"city"`
	District     sql.NullString `json:"district"`
	Phone        string         `json:"phone"`
	Whatsapp     string         `json:"whatsapp"`
	PhotoUrl     []byte         `json:"photo_url"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.PasswordHash,
		arg.Role,
		arg.Country,
		arg.City,
		arg.District,
		arg.Phone,
		arg.Whatsapp,
		arg.PhotoUrl,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Country,
		&i.City,
		&i.District,
		&i.Phone,
		&i.Whatsapp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhotoUrl,
		&i.Description,
		&i.IsVerified,
		&i.IsBlocked,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getBlockerUser = `-- name: GetBlockerUser :one
SELECT id,
    is_blocked
FROM users
WHERE id = $1
`

type GetBlockerUserRow struct {
	ID        int64        `json:"id"`
	IsBlocked sql.NullBool `json:"is_blocked"`
}

func (q *Queries) GetBlockerUser(ctx context.Context, id int64) (GetBlockerUserRow, error) {
	row := q.db.QueryRowContext(ctx, getBlockerUser, id)
	var i GetBlockerUserRow
	err := row.Scan(&i.ID, &i.IsBlocked)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Country,
		&i.City,
		&i.District,
		&i.Phone,
		&i.Whatsapp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhotoUrl,
		&i.Description,
		&i.IsVerified,
		&i.IsBlocked,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Country,
		&i.City,
		&i.District,
		&i.Phone,
		&i.Whatsapp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhotoUrl,
		&i.Description,
		&i.IsVerified,
		&i.IsBlocked,
	)
	return i, err
}

const getUserByIDFromAdmin = `-- name: GetUserByIDFromAdmin :one
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByIDFromAdmin(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByIDFromAdmin, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Country,
		&i.City,
		&i.District,
		&i.Phone,
		&i.Whatsapp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhotoUrl,
		&i.Description,
		&i.IsVerified,
		&i.IsBlocked,
	)
	return i, err
}

const getUserByIDFromUser = `-- name: GetUserByIDFromUser :one
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByIDFromUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByIDFromUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Country,
		&i.City,
		&i.District,
		&i.Phone,
		&i.Whatsapp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhotoUrl,
		&i.Description,
		&i.IsVerified,
		&i.IsBlocked,
	)
	return i, err
}

const listBlockedUsers = `-- name: ListBlockedUsers :many
SELECT id,
    is_blocked
FROM users
WHERE is_blocked = true
`

type ListBlockedUsersRow struct {
	ID        int64        `json:"id"`
	IsBlocked sql.NullBool `json:"is_blocked"`
}

func (q *Queries) ListBlockedUsers(ctx context.Context) ([]ListBlockedUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listBlockedUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListBlockedUsersRow{}
	for rows.Next() {
		var i ListBlockedUsersRow
		if err := rows.Scan(&i.ID, &i.IsBlocked); err != nil {
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

const listPartners = `-- name: ListPartners :many
SELECT id,
    username,
    email
FROM users
WHERE role = 'partner'
`

type ListPartnersRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (q *Queries) ListPartners(ctx context.Context) ([]ListPartnersRow, error) {
	rows, err := q.db.QueryContext(ctx, listPartners)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListPartnersRow{}
	for rows.Next() {
		var i ListPartnersRow
		if err := rows.Scan(&i.ID, &i.Username, &i.Email); err != nil {
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

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.PasswordHash,
			&i.PasswordChangeAt,
			&i.Role,
			&i.Country,
			&i.City,
			&i.District,
			&i.Phone,
			&i.Whatsapp,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PhotoUrl,
			&i.Description,
			&i.IsVerified,
			&i.IsBlocked,
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

const listUsersByEmail = `-- name: ListUsersByEmail :many
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
WHERE email ILIKE '%' || $1 || '%'
ORDER BY email
`

func (q *Queries) ListUsersByEmail(ctx context.Context, dollar_1 sql.NullString) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsersByEmail, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.PasswordHash,
			&i.PasswordChangeAt,
			&i.Role,
			&i.Country,
			&i.City,
			&i.District,
			&i.Phone,
			&i.Whatsapp,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PhotoUrl,
			&i.Description,
			&i.IsVerified,
			&i.IsBlocked,
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

const listUsersByRole = `-- name: ListUsersByRole :many
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
WHERE role = $1
ORDER BY username
`

func (q *Queries) ListUsersByRole(ctx context.Context, role NullRole) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsersByRole, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.PasswordHash,
			&i.PasswordChangeAt,
			&i.Role,
			&i.Country,
			&i.City,
			&i.District,
			&i.Phone,
			&i.Whatsapp,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PhotoUrl,
			&i.Description,
			&i.IsVerified,
			&i.IsBlocked,
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

const listUsersByUsername = `-- name: ListUsersByUsername :many
SELECT id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
FROM users
WHERE username ILIKE '%' || $1 || '%'
ORDER BY username
`

func (q *Queries) ListUsersByUsername(ctx context.Context, dollar_1 sql.NullString) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsersByUsername, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.PasswordHash,
			&i.PasswordChangeAt,
			&i.Role,
			&i.Country,
			&i.City,
			&i.District,
			&i.Phone,
			&i.Whatsapp,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PhotoUrl,
			&i.Description,
			&i.IsVerified,
			&i.IsBlocked,
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

const unblockUser = `-- name: UnblockUser :one
UPDATE users
SET is_blocked = false
WHERE id = $1
RETURNING is_blocked
`

func (q *Queries) UnblockUser(ctx context.Context, id int64) (sql.NullBool, error) {
	row := q.db.QueryRowContext(ctx, unblockUser, id)
	var is_blocked sql.NullBool
	err := row.Scan(&is_blocked)
	return is_blocked, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET username = $2,
    email = $3,
    country = $4,
    city = $5,
    district = $6,
    phone = $7,
    whatsapp = $8,
    photo_url = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
`

type UpdateUserParams struct {
	ID       int64          `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Country  sql.NullString `json:"country"`
	City     sql.NullString `json:"city"`
	District sql.NullString `json:"district"`
	Phone    string         `json:"phone"`
	Whatsapp string         `json:"whatsapp"`
	PhotoUrl []byte         `json:"photo_url"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.Country,
		arg.City,
		arg.District,
		arg.Phone,
		arg.Whatsapp,
		arg.PhotoUrl,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Country,
		&i.City,
		&i.District,
		&i.Phone,
		&i.Whatsapp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhotoUrl,
		&i.Description,
		&i.IsVerified,
		&i.IsBlocked,
	)
	return i, err
}

const updateUserForAdmin = `-- name: UpdateUserForAdmin :one
UPDATE users
SET username = $2,
    email = $3,
    country = $4,
    city = $5,
    district = $6,
    phone = $7,
    whatsapp = $8,
    photo_url = $9,
    role = $10,
    updated_at = NOW()
WHERE id = $1
RETURNING id, username, email, password_hash, password_change_at, role, country, city, district, phone, whatsapp, created_at, updated_at, photo_url, description, is_verified, is_blocked
`

type UpdateUserForAdminParams struct {
	ID       int64          `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Country  sql.NullString `json:"country"`
	City     sql.NullString `json:"city"`
	District sql.NullString `json:"district"`
	Phone    string         `json:"phone"`
	Whatsapp string         `json:"whatsapp"`
	PhotoUrl []byte         `json:"photo_url"`
	Role     NullRole       `json:"role"`
}

func (q *Queries) UpdateUserForAdmin(ctx context.Context, arg UpdateUserForAdminParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserForAdmin,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.Country,
		arg.City,
		arg.District,
		arg.Phone,
		arg.Whatsapp,
		arg.PhotoUrl,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Country,
		&i.City,
		&i.District,
		&i.Phone,
		&i.Whatsapp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhotoUrl,
		&i.Description,
		&i.IsVerified,
		&i.IsBlocked,
	)
	return i, err
}
