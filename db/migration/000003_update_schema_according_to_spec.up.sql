-- 1. Добавление полей в таблицу users для более детального профиля
ALTER TABLE "users"
ADD COLUMN IF NOT EXISTS "photo_url" varchar;
ALTER TABLE "users"
ADD COLUMN IF NOT EXISTS "description" text;
ALTER TABLE "users"
ADD COLUMN IF NOT EXISTS "is_verified" boolean DEFAULT false;
ALTER TABLE "users"
ADD COLUMN IF NOT EXISTS "is_blocked" boolean DEFAULT false;
-- 2. Добавление новых категорий услуг согласно ТЗ
INSERT INTO "service_categories" ("name")
VALUES ('Аренда авто с водителем'),
    ('Аренда спецтехники'),
    ('Грузовые перевозки'),
    ('Караоке такси') ON CONFLICT ("name") DO NOTHING;
-- 3. Добавление поля "subcategory" в таблицу services для более детальной классификации
ALTER TABLE "services"
ADD COLUMN IF NOT EXISTS "subcategory" varchar;
-- 4. Добавление полей для услуг с возможностью указать город и район
ALTER TABLE "services"
ADD COLUMN IF NOT EXISTS "country" varchar;
ALTER TABLE "services"
ADD COLUMN IF NOT EXISTS "city" varchar;
ALTER TABLE "services"
ADD COLUMN IF NOT EXISTS "district" varchar;
-- 5. Создание таблицы для чата между клиентами и услугодателями
CREATE TABLE IF NOT EXISTS "messages" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "sender_id" bigint NOT NULL REFERENCES "users" ("id"),
    "receiver_id" bigint NOT NULL REFERENCES "users" ("id"),
    "content" text NOT NULL,
    "is_read" boolean DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
-- Индексы для быстрого поиска сообщений
CREATE INDEX IF NOT EXISTS "messages_sender_id_idx" ON "messages" ("sender_id");
CREATE INDEX IF NOT EXISTS "messages_receiver_id_idx" ON "messages" ("receiver_id");
CREATE INDEX IF NOT EXISTS "messages_created_at_idx" ON "messages" ("created_at");
-- 6. Расширение таблицы promo_codes для привязки к партнерам
ALTER TABLE "promo_codes"
ADD COLUMN IF NOT EXISTS "max_usages" integer DEFAULT 1;
ALTER TABLE "promo_codes"
ADD COLUMN IF NOT EXISTS "current_usages" integer DEFAULT 0;
-- 7. Расширение таблицы subscriptions для дополнительной информации
ALTER TABLE "subscriptions"
ADD COLUMN IF NOT EXISTS "subscription_type" varchar;
ALTER TABLE "subscriptions"
ADD COLUMN IF NOT EXISTS "price" decimal(10, 2);
ALTER TABLE "subscriptions"
ADD COLUMN IF NOT EXISTS "promo_code_id" bigint REFERENCES "promo_codes" ("id");
-- 8. Улучшения для таблицы orders
ALTER TABLE "orders"
ADD COLUMN IF NOT EXISTS "provider_accepted" boolean DEFAULT false;
ALTER TABLE "orders"
ADD COLUMN IF NOT EXISTS "provider_message" text;
ALTER TABLE "orders"
ADD COLUMN IF NOT EXISTS "client_message" text;
ALTER TABLE "orders"
ADD COLUMN IF NOT EXISTS "order_date" timestamptz;
-- 9. Создание таблицы для избранных услугодателей клиента
CREATE TABLE IF NOT EXISTS "favorites" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "client_id" bigint NOT NULL REFERENCES "users" ("id"),
    "provider_id" bigint NOT NULL REFERENCES "users" ("id"),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    UNIQUE ("client_id", "provider_id")
);
-- 10. Создание таблицы для статистики партнеров
CREATE TABLE IF NOT EXISTS "partner_statistics" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "partner_id" bigint NOT NULL REFERENCES "users" ("id"),
    "providers_attracted" integer DEFAULT 0,
    "total_subscriptions" integer DEFAULT 0,
    "active_subscriptions" integer DEFAULT 0,
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);
-- 11. Добавление индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS "services_city_district_idx" ON "services" ("city", "district");
CREATE INDEX IF NOT EXISTS "services_price_rating_idx" ON "services" ("price");
CREATE INDEX IF NOT EXISTS "orders_provider_accepted_idx" ON "orders" ("provider_accepted");
CREATE INDEX IF NOT EXISTS "orders_order_date_idx" ON "orders" ("order_date");