CREATE TABLE IF NOT EXISTS "city" (
    "id" bigserial PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE
);