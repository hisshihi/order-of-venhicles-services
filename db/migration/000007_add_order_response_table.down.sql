-- Удаление индекса для поиска по выбранному провайдеру
DROP INDEX IF EXISTS "orders_selected_provider_id_idx";

-- Удаление ссылки на выбранного провайдера
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_selected_provider_id_fkey";
ALTER TABLE "orders" DROP COLUMN IF EXISTS "selected_provider_id";

-- Удаление индексов таблицы откликов
DROP INDEX IF EXISTS "order_responses_order_provider_unique_idx";
DROP INDEX IF EXISTS "order_responses_is_selected_idx";
DROP INDEX IF EXISTS "order_responses_provider_id_idx";
DROP INDEX IF EXISTS "order_responses_order_id_idx";

-- Удаление таблицы откликов
DROP TABLE IF EXISTS "order_responses";