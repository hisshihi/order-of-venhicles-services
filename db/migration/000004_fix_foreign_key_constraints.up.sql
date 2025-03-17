-- Удаляем неправильное ограничение внешнего ключа
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_category_id_fkey";
-- Изменяем ссылку поля category_id на правильную таблицу service_categories
ALTER TABLE "orders"
ADD CONSTRAINT "orders_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "service_categories" ("id");
-- Комментарий: Исправлено ограничение внешнего ключа для поля category_id
-- Теперь поле ссылается на таблицу service_categories, а не на services