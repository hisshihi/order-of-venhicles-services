-- Восстанавливаем уникальные ограничения
ALTER TABLE users 
  ADD CONSTRAINT users_phone_key UNIQUE (phone);

ALTER TABLE users 
  ADD CONSTRAINT users_whatsapp_key UNIQUE (whatsapp);