ALTER TABLE "reviews" DROP CONSTRAINT IF EXISTS unique_provider_id;
ALTER TABLE "reviews"
ADD CONSTRAINT unique_provider_client UNIQUE (provider_id, client_id);