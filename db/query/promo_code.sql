-- name: CreatePromoCode :one
INSERT INTO promo_codes (partner_id, code, discount_percentage, valid_until, max_usages, current_usages)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPromoCodeByID :one
SELECT * FROM promo_codes
WHERE id = $1;

-- name: GetPromoCodeByCode :one
SELECT * FROM promo_codes
WHERE code = $1;

-- name: GetPromoCodeByPartnerID :one
SELECT * FROM promo_codes
WHERE partner_id = $1;

-- name: ListPromoCodes :many
SELECT * FROM promo_codes
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPromoCodesByPartnerID :many
SELECT * FROM promo_codes
WHERE partner_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProvidersByPromoCode :many
SELECT DISTINCT u.id, u.username, u.email, u.phone, u.city, s.id as subscription_id, 
       s.start_date, s.end_date, s.status, s.subscription_type, 
       s.price, s.original_price, p.discount_percentage
FROM users u
JOIN subscriptions s ON u.id = s.provider_id
JOIN promo_codes p ON s.promo_code_id = p.id
WHERE p.partner_id = $1 AND p.id = $2
ORDER BY s.status DESC
LIMIT $3 OFFSET $4;

-- name: GetAllProvidersByPartnerPromos :many
-- Получает всех поставщиков, которые использовали данный промокод
SELECT DISTINCT u.id, u.username, u.email, u.phone, u.city, s.id as subscription_id, 
       s.start_date, s.end_date, s.status, s.subscription_type, 
       s.price, s.original_price, p.code, p.discount_percentage
FROM users u
JOIN subscriptions s ON u.id = s.provider_id
JOIN promo_codes p ON s.promo_code_id = p.id
WHERE p.partner_id = $1
ORDER BY s.status DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePromoCodeByID :exec
UPDATE promo_codes
SET current_usages = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeletePromoCode :exec
DELETE FROM promo_codes
WHERE id = $1;

-- name: CountPromoCode :one
SELECT COUNT(*) FROM "promo_codes";