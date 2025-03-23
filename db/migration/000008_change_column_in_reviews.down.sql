-- 1. Удаляем уникальное ограничение на provider_id
ALTER TABLE reviews DROP CONSTRAINT unique_provider_id;
-- 2. Возвращаем колонку order_id
ALTER TABLE reviews
ADD COLUMN order_id bigint;
-- 3. Устанавливаем NOT NULL для order_id и добавляем внешний ключ
ALTER TABLE reviews
ALTER COLUMN order_id
SET NOT NULL,
    ADD CONSTRAINT reviews_order_id_fkey FOREIGN KEY (order_id) REFERENCES orders(id);
-- 4. Убираем NOT NULL с provider_id (если это было изначально)
ALTER TABLE reviews
ALTER COLUMN provider_id DROP NOT NULL;