CREATE TYPE "role" AS ENUM (
    'provider',
    'client',
    'partner',
    'admin'
);
CREATE TYPE "status_orders" AS ENUM (
    'pending',
    'accepted',
    'completed',
    'cancelled'
);
CREATE TYPE "status_subscription" AS ENUM ('active', 'inactive', 'expired');
CREATE TYPE "status_payment" AS ENUM ('pending', 'completed', 'failed');
CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "username" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password_hash" varchar NOT NULL,
    "password_change_at" timestamptz NOT NULL DEFAULT null,
    "role" role,
    "country" varchar,
    "city" varchar,
    "district" varchar,
    "phone" varchar UNIQUE NOT NULL,
    "whatsapp" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "service_categories" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "name" varchar NOT NULL,
    "description" text NOT NULL
);
CREATE TABLE "services" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "provider_id" bigint NOT NULL,
    "category_id" bigint NOT NULL,
    "title" varchar NOT NULL,
    "description" text NOT NULL,
    "price" decimal(10, 2) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "orders" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "client_id" bigint NOT NULL,
    "service_id" bigint NOT NULL,
    "status" "status_orders",
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "subscriptions" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "provider_id" bigint NOT NULL,
    "start_date" date NOT NULL,
    "end_date" date NOT NULL,
    "status" "status_subscription",
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "promo_codes" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "partner_id" bigint NOT NULL,
    "code" varchar UNIQUE NOT NULL,
    "discount_percentage" integer NOT NULL,
    "valid_until" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "reviews" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "order_id" bigint UNIQUE NOT NULL,
    "client_id" bigint NOT NULL,
    "provider_id" bigint NOT NULL,
    "rating" integer NOT NULL CHECK (
        rating >= 1
        AND rating <= 5
    ),
    "comment" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "payments" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "user_id" bigint NOT NULL,
    "amount" decimal(10, 2) NOT NULL,
    "payment_method" varchar NOT NULL,
    "status" "status_payment",
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE INDEX ON "users" ("role");
CREATE INDEX ON "users" ("country", "city", "district");
CREATE INDEX ON "service_categories" ("name");
CREATE INDEX ON "services" ("provider_id");
CREATE INDEX ON "services" ("category_id");
CREATE INDEX ON "services" ("price");
CREATE INDEX ON "orders" ("client_id");
CREATE INDEX ON "orders" ("service_id");
CREATE INDEX ON "orders" ("status");
CREATE INDEX ON "subscriptions" ("provider_id");
CREATE INDEX ON "subscriptions" ("status");
CREATE INDEX ON "promo_codes" ("partner_id");
CREATE INDEX ON "promo_codes" ("code");
CREATE INDEX ON "reviews" ("order_id");
CREATE INDEX ON "reviews" ("client_id");
CREATE INDEX ON "reviews" ("provider_id");
CREATE INDEX ON "payments" ("user_id");
CREATE INDEX ON "payments" ("status");
ALTER TABLE "services"
ADD FOREIGN KEY ("provider_id") REFERENCES "users" ("id");
ALTER TABLE "services"
ADD FOREIGN KEY ("category_id") REFERENCES "service_categories" ("id");
ALTER TABLE "orders"
ADD FOREIGN KEY ("client_id") REFERENCES "users" ("id");
ALTER TABLE "orders"
ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE "subscriptions"
ADD FOREIGN KEY ("provider_id") REFERENCES "users" ("id");
ALTER TABLE "promo_codes"
ADD FOREIGN KEY ("partner_id") REFERENCES "users" ("id");
ALTER TABLE "reviews"
ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "reviews"
ADD FOREIGN KEY ("client_id") REFERENCES "users" ("id");
ALTER TABLE "reviews"
ADD FOREIGN KEY ("provider_id") REFERENCES "users" ("id");
ALTER TABLE "payments"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");