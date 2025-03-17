-- name: CreateMessage :one
-- Создает новое сообщение
INSERT INTO "messages" (sender_id, receiver_id, content)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetMessagesByUsers :many
-- Получает историю переписки между двумя пользователями
SELECT m.*,
    s.username as sender_name,
    r.username as receiver_name
FROM "messages" m
    JOIN "users" s ON m.sender_id = s.id
    JOIN "users" r ON m.receiver_id = r.id
WHERE (
        m.sender_id = $1
        AND m.receiver_id = $2
    )
    OR (
        m.sender_id = $2
        AND m.receiver_id = $1
    )
ORDER BY m.created_at ASC
LIMIT $3 OFFSET $4;
-- name: GetUnreadMessagesCount :one
-- Получает количество непрочитанных сообщений для пользователя
SELECT COUNT(*)
FROM "messages"
WHERE receiver_id = $1
    AND is_read = false;
-- name: MarkMessagesAsRead :exec
-- Отмечает сообщения как прочитанные
UPDATE "messages"
SET is_read = true
WHERE receiver_id = $1
    AND sender_id = $2
    AND is_read = false;
-- name: GetUserRecentChats :many
-- Получает список недавних чатов пользователя
WITH recent_messages AS (
    SELECT DISTINCT ON (
            CASE
                WHEN sender_id = $1 THEN receiver_id
                ELSE sender_id
            END
        ) id,
        CASE
            WHEN sender_id = $1 THEN receiver_id
            ELSE sender_id
        END as other_user_id,
        content,
        created_at
    FROM "messages"
    WHERE sender_id = $1
        OR receiver_id = $1
    ORDER BY other_user_id,
        created_at DESC
)
SELECT rm.id as message_id,
    rm.other_user_id,
    u.username as other_user_name,
    u.photo_url as other_user_photo,
    rm.content as last_message,
    rm.created_at as message_time,
    (
        SELECT COUNT(*)
        FROM "messages"
        WHERE sender_id = rm.other_user_id
            AND receiver_id = $1
            AND is_read = false
    ) as unread_count
FROM recent_messages rm
    JOIN "users" u ON rm.other_user_id = u.id
ORDER BY rm.created_at DESC
LIMIT $2 OFFSET $3;