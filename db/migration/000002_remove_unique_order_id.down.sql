-- Удаляем новые индексы
DROP INDEX IF EXISTS "reviews_client_id_order_id_key";
-- Восстанавливаем оригинальное ограничение UNIQUE на order_id
ALTER TABLE "reviews"
ADD CONSTRAINT "reviews_order_id_key" UNIQUE ("order_id");