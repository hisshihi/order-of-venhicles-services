-- Удаляем исправленное ограничение внешнего ключа
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_category_id_fkey";
-- Возвращаем ссылку поля category_id на таблицу services (как было раньше)
ALTER TABLE "orders"
ADD CONSTRAINT "orders_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "services" ("id");
-- Комментарий: Восстановлено исходное ограничение внешнего ключа для поля category_id