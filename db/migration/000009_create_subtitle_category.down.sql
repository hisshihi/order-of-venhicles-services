-- Удаление внешнего ключа из таблицы services
ALTER TABLE "services" DROP CONSTRAINT fk_services_subtitle_category;
-- Удаление колонки subtitle_category_id из таблицы services
ALTER TABLE "services" DROP COLUMN "subtitle_category_id";
-- Удаление внешнего ключа из таблицы orders
ALTER TABLE "orders" DROP CONSTRAINT fk_orders_subtitle_category;
-- Удаление колонки subtitle_category_id из таблицы orders
ALTER TABLE "orders" DROP COLUMN "subtitle_category_id";
-- Удаление индекса
DROP INDEX idx_subtitle_category_name;
-- Удаление таблицы subtitle_category
DROP TABLE "subtitle_category";