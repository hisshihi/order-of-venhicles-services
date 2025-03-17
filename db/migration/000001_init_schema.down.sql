-- Сначала удаляем внешние ключи
ALTER TABLE "services" DROP CONSTRAINT IF EXISTS "services_provider_id_fkey";
ALTER TABLE "services" DROP CONSTRAINT IF EXISTS "services_category_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_client_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_category_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_service_id_fkey";
ALTER TABLE "subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_provider_id_fkey";
ALTER TABLE "promo_codes" DROP CONSTRAINT IF EXISTS "promo_codes_partner_id_fkey";
ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS "reviews_order_id_fkey";
ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS "reviews_client_id_fkey";
ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS "reviews_provider_id_fkey";
ALTER TABLE "payments" DROP CONSTRAINT IF EXISTS "payments_user_id_fkey";
-- Затем удаляем таблицы
DROP TABLE IF EXISTS "payments";
DROP TABLE IF EXISTS "reviews";
DROP TABLE IF EXISTS "promo_codes";
DROP TABLE IF EXISTS "subscriptions";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "services";
DROP TABLE IF EXISTS "service_categories";
DROP TABLE IF EXISTS "users";
-- В конце удаляем типы
DROP TYPE IF EXISTS "status_payment";
DROP TYPE IF EXISTS "status_subscription";
DROP TYPE IF EXISTS "status_orders";
DROP TYPE IF EXISTS "role";