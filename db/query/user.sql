-- name: CreateUser :one
INSERT INTO users (
        username,
        email,
        password_hash,
        role,
        country,
        city,
        district,
        phone,
        whatsapp
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
        $9
    )
RETURNING *;

-- name: GetUserByIDFromAdmin :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByIDFromUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: ListUsersByUsername :many
SELECT * FROM users
WHERE username ILIKE '%' || $1 || '%'
ORDER BY username;

-- name: ListUsersByEmail :many
SELECT * FROM users
WHERE email ILIKE '%' || $1 || '%'
ORDER BY email;

-- name: ListUsersByRole :many
SELECT * FROM users
WHERE role = $1
ORDER BY username;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET username = $2,
    email = $3,
    country = $4,
    city = $5,
    district = $6,
    phone = $7,
    whatsapp = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;