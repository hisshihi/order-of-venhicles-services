-- Удаление индексов для оптимизации запросов
DROP INDEX IF EXISTS "orders_order_date_idx";
DROP INDEX IF EXISTS "orders_provider_accepted_idx";
DROP INDEX IF EXISTS "services_price_rating_idx";
DROP INDEX IF EXISTS "services_city_district_idx";
-- Удаление таблиц
DROP TABLE IF EXISTS "partner_statistics";
DROP TABLE IF EXISTS "favorites";
DROP TABLE IF EXISTS "messages";
-- Удаление полей из таблицы orders
ALTER TABLE "orders" DROP COLUMN IF EXISTS "order_date";
ALTER TABLE "orders" DROP COLUMN IF EXISTS "client_message";
ALTER TABLE "orders" DROP COLUMN IF EXISTS "provider_message";
ALTER TABLE "orders" DROP COLUMN IF EXISTS "provider_accepted";
-- Удаление полей из таблицы subscriptions
ALTER TABLE "subscriptions" DROP COLUMN IF EXISTS "promo_code_id";
ALTER TABLE "subscriptions" DROP COLUMN IF EXISTS "price";
ALTER TABLE "subscriptions" DROP COLUMN IF EXISTS "subscription_type";
-- Удаление полей из таблицы promo_codes
ALTER TABLE "promo_codes" DROP COLUMN IF EXISTS "current_usages";
ALTER TABLE "promo_codes" DROP COLUMN IF EXISTS "max_usages";
-- Удаление полей из таблицы services
ALTER TABLE "services" DROP COLUMN IF EXISTS "district";
ALTER TABLE "services" DROP COLUMN IF EXISTS "city";
ALTER TABLE "services" DROP COLUMN IF EXISTS "country";
ALTER TABLE "services" DROP COLUMN IF EXISTS "subcategory";
-- Удаление полей из таблицы users
ALTER TABLE "users" DROP COLUMN IF EXISTS "is_blocked";
ALTER TABLE "users" DROP COLUMN IF EXISTS "is_verified";
ALTER TABLE "users" DROP COLUMN IF EXISTS "description";
ALTER TABLE "users" DROP COLUMN IF EXISTS "photo_url";