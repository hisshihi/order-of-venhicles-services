DROP INDEX IF EXISTS "subscriptions_promo_code_id_idx";

ALTER TABLE "subscriptions" DROP CONSTRAINT IF EXISTS "subscriptions_promo_code_id_fkey";

ALTER TABLE "subscriptions" DROP COLUMN IF EXISTS "promo_code_id";
ALTER TABLE "subscriptions" DROP COLUMN IF EXISTS "price";
ALTER TABLE "subscriptions" DROP COLUMN IF EXISTS "original_price";
ALTER TABLE "subscriptions" DROP COLUMN IF EXISTS "subscription_type";