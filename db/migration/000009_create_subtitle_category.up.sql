-- Создание таблицы subtitle_category
CREATE TABLE "subtitle_category" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
-- Создание индекса на поле name для ускорения поиска
CREATE INDEX idx_subtitle_category_name ON "subtitle_category" ("name");
-- Добавление колонки subtitle_category_id в таблицу services
ALTER TABLE "services"
ADD COLUMN "subtitle_category_id" bigint;
-- Добавление внешнего ключа в таблицу services
ALTER TABLE "services"
ADD CONSTRAINT fk_services_subtitle_category FOREIGN KEY ("subtitle_category_id") REFERENCES "subtitle_category" ("id");
-- Добавление колонки subtitle_category_id в таблицу orders
ALTER TABLE "orders"
ADD COLUMN "subtitle_category_id" bigint;
-- Добавление внешнего ключа в таблицу orders
ALTER TABLE "orders"
ADD CONSTRAINT fk_orders_subtitle_category FOREIGN KEY ("subtitle_category_id") REFERENCES "subtitle_category" ("id");