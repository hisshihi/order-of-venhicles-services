-- Удаляем внешний ключ
ALTER TABLE support_messages
DROP CONSTRAINT fk_sender;

-- Удаляем таблицу
DROP TABLE IF EXISTS support_messages;