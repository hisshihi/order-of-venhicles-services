-- name: CreateSupportMessage :one
INSERT INTO "support_messages" ("sender_id", "subject", "messages")
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListSupportMessages :many
SELECT * FROM "support_messages"
ORDER BY "created_at";

-- name: CountSupportMessages :one
SELECT COUNT(*) FROM "support_messages";