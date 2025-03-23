-- 1. Удаляем внешний ключ, связанный с order_id
ALTER TABLE reviews DROP CONSTRAINT reviews_order_id_fkey;
-- 2. Удаляем колонку order_id
ALTER TABLE reviews DROP COLUMN order_id;
-- 3. Изменяем provider_id: тип bigint, NOT NULL, уникальное ограничение
ALTER TABLE reviews
ALTER COLUMN provider_id TYPE bigint,
    ALTER COLUMN provider_id
SET NOT NULL,
    ADD CONSTRAINT unique_provider_id UNIQUE (provider_id);