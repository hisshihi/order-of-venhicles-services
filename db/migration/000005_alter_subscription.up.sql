ALTER TABLE "subscriptions" ADD COLUMN IF NOT EXISTS "promo_code_id" bigint REFERENCES "promo_codes" ("id");
ALTER TABLE "subscriptions" ADD COLUMN IF NOT EXISTS "price" decimal(10, 2);
ALTER TABLE "subscriptions" ADD COLUMN IF NOT EXISTS "original_price" decimal(10, 2);
ALTER TABLE "subscriptions" ADD COLUMN IF NOT EXISTS "subscription_type" varchar;

CREATE INDEX IF NOT EXISTS "subscriptions_promo_code_id_idx" ON "subscriptions" ("promo_code_id");