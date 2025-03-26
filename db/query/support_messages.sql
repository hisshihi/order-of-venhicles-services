-- name: CreateSupportMessage :one
INSERT INTO "support_messages" ("sender_id", "subject", "messages")
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListSupportMessages :many
SELECT * FROM "support_messages"
ORDER BY "created_at"
LIMIT $1 OFFSET $2;

-- name: CountSupportMessages :one
SELECT COUNT(*) FROM "support_messages";

-- name: DeleteSupportMessage :exec
DELETE FROM "support_messages"
WHERE id = $1;