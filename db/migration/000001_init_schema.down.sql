-- Сначала удаляем внешние ключи
ALTER TABLE "Service" DROP CONSTRAINT IF EXISTS "Service_provider_id_fkey";
ALTER TABLE "Service" DROP CONSTRAINT IF EXISTS "Service_category_id_fkey";
ALTER TABLE "Orders" DROP CONSTRAINT IF EXISTS "Orders_client_id_fkey";
ALTER TABLE "Orders" DROP CONSTRAINT IF EXISTS "Orders_service_id_fkey";
ALTER TABLE "Subscriptions" DROP CONSTRAINT IF EXISTS "Subscriptions_provider_id_fkey";
ALTER TABLE "PromoCodes" DROP CONSTRAINT IF EXISTS "PromoCodes_partner_id_fkey";
ALTER TABLE "Reviews" DROP CONSTRAINT IF EXISTS "Reviews_order_id_fkey";
ALTER TABLE "Reviews" DROP CONSTRAINT IF EXISTS "Reviews_client_id_fkey";
ALTER TABLE "Reviews" DROP CONSTRAINT IF EXISTS "Reviews_provider_id_fkey";
ALTER TABLE "Payment" DROP CONSTRAINT IF EXISTS "Payment_user_id_fkey";
-- Затем удаляем таблицы
DROP TABLE IF EXISTS "Payment";
DROP TABLE IF EXISTS "Reviews";
DROP TABLE IF EXISTS "PromoCodes";
DROP TABLE IF EXISTS "Subscriptions";
DROP TABLE IF EXISTS "Orders";
DROP TABLE IF EXISTS "Service";
DROP TABLE IF EXISTS "ServiceCategories";
DROP TABLE IF EXISTS "Users";
-- В конце удаляем типы
DROP TYPE IF EXISTS "status_payment";
DROP TYPE IF EXISTS "status_subscription";
DROP TYPE IF EXISTS "status_orders";
DROP TYPE IF EXISTS "role";