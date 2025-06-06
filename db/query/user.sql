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
RETURNING *;
-- name: GetUserByIDFromAdmin :one
SELECT *
FROM users
WHERE id = $1;
-- name: GetUserByIDFromUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;
-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;
-- name: ListUsersByUsername :many
SELECT *
FROM users
WHERE username ILIKE '%' || $1 || '%'
ORDER BY username;
-- name: ListUsersByEmail :many
SELECT *
FROM users
WHERE email ILIKE '%' || $1 || '%'
ORDER BY email;
-- name: ListUsersByRole :many
SELECT *
FROM users
WHERE role = $1
ORDER BY username;
-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
-- name: CountUsers :one
SELECT COUNT(*)
FROM users;
-- name: UpdateUser :one
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
RETURNING *;
-- name: UpdateUserForAdmin :one
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
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
-- name: ChangePassword :exec
UPDATE users
SET password_hash = $2,
    password_change_at = NOW()
WHERE id = $1;
-- name: ListPartners :many
SELECT id,
    username,
    email
FROM users
WHERE role = 'partner';
-- name: BlockedUser :one
UPDATE users
SET is_blocked = true
WHERE id = $1
RETURNING is_blocked;
-- name: ListBlockedUsers :many
SELECT id,
    is_blocked
FROM users
WHERE is_blocked = true;
-- name: GetBlockerUser :one
SELECT id,
    is_blocked
FROM users
WHERE id = $1;
-- name: UnblockUser :one
UPDATE users
SET is_blocked = false
WHERE id = $1
RETURNING is_blocked;