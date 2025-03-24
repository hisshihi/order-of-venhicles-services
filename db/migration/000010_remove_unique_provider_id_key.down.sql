ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS unique_provider_client;
ALTER TABLE "reviews"
ADD CONSTRAINT unique_provider_id UNIQUE (provider_id);