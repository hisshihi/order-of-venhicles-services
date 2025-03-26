CREATE TABLE support_messages (
    "id" bigserial PRIMARY KEY,
    "sender_id" bigint NOT NULL,
    "subject" VARCHAR(255) NOT NULL,
    "messages" TEXT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE support_messages
ADD CONSTRAINT fk_sender
FOREIGN KEY (sender_id)
REFERENCES users(id)
ON DELETE CASCADE;