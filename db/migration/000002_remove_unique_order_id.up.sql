-- Удаляем существующее ограничение UNIQUE с order_id
ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS "reviews_order_id_key";
-- Создаем составной уникальный индекс на (client_id, order_id)
-- Это позволит разным клиентам оставлять отзывы на один заказ,
-- но один клиент сможет оставить только один отзыв на конкретный заказ
CREATE UNIQUE INDEX IF NOT EXISTS "reviews_client_id_order_id_key" ON "reviews" ("client_id", "order_id");