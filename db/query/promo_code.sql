-- name: CreatePromoCode :one
INSERT INTO promo_codes (partner_id, code, discount_percentage, valid_until)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPromoCodeByID :one
SELECT * FROM promo_codes
WHERE id = $1;

-- name: GetPromoCodeByPartnerID :one
SELECT * FROM promo_codes
WHERE partner_id = $1;

-- name: ListPromoCodes :many
SELECT * FROM promo_codes
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: UpdatePromoCode :one
UPDATE promo_codes
SET partner_id = $2, code = $3, discount_percentage = $4, valid_until = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePromoCode :exec
DELETE FROM promo_codes
WHERE id = $1;