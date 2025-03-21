-- Создание таблицы для откликов услугодателей на заказы
CREATE TABLE IF NOT EXISTS "order_responses" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "order_id" bigint NOT NULL REFERENCES "orders" ("id"),
    "provider_id" bigint NOT NULL REFERENCES "users" ("id"),
    "message" text,
    "offered_price" decimal(10, 2),
    "is_selected" boolean DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS "order_responses_order_id_idx" ON "order_responses" ("order_id");
CREATE INDEX IF NOT EXISTS "order_responses_provider_id_idx" ON "order_responses" ("provider_id");
CREATE INDEX IF NOT EXISTS "order_responses_is_selected_idx" ON "order_responses" ("is_selected");

-- Уникальное ограничение - один провайдер может оставить только один отклик на заказ
CREATE UNIQUE INDEX IF NOT EXISTS "order_responses_order_provider_unique_idx" ON "order_responses" ("order_id", "provider_id");

-- Добавим поле selected_provider_id, которое будет указывать на выбранного провайдера
ALTER TABLE "orders" ADD COLUMN IF NOT EXISTS "selected_provider_id" bigint REFERENCES "users" ("id");

-- Индекс для быстрого поиска по выбранному провайдеру
CREATE INDEX IF NOT EXISTS "orders_selected_provider_id_idx" ON "orders" ("selected_provider_id");