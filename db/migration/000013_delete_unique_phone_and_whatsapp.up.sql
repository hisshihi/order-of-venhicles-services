-- Удаляем уникальные ограничения с phone и whatsapp
ALTER TABLE users 
  DROP CONSTRAINT IF EXISTS users_phone_key;

ALTER TABLE users 
  DROP CONSTRAINT IF EXISTS users_whatsapp_key;