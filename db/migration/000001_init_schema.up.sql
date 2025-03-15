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
CREATE TABLE "Users" (
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
CREATE TABLE "ServiceCategories" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "name" varchar NOT NULL,
    "description" text NOT NULL
);
CREATE TABLE "Service" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "provider_id" bigint NOT NULL,
    "category_id" bigint NOT NULL,
    "title" varchar NOT NULL,
    "description" text NOT NULL,
    "price" decimal(10, 2) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "Orders" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "client_id" bigint NOT NULL,
    "service_id" bigint NOT NULL,
    "status" "status_orders",
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "Subscriptions" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "provider_id" bigint NOT NULL,
    "start_date" date NOT NULL,
    "end_date" date NOT NULL,
    "status" "status_subscription",
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "PromoCodes" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "partner_id" bigint NOT NULL,
    "code" varchar UNIQUE NOT NULL,
    "discount_percentage" integer NOT NULL,
    "valid_until" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE TABLE "Reviews" (
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
CREATE TABLE "Payment" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "user_id" bigint NOT NULL,
    "amount" decimal(10, 2) NOT NULL,
    "payment_method" varchar NOT NULL,
    "status" "status_payment",
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
CREATE INDEX ON "Users" ("role");
CREATE INDEX ON "Users" ("country", "city", "district");
CREATE INDEX ON "ServiceCategories" ("name");
CREATE INDEX ON "Service" ("provider_id");
CREATE INDEX ON "Service" ("category_id");
CREATE INDEX ON "Service" ("price");
CREATE INDEX ON "Orders" ("client_id");
CREATE INDEX ON "Orders" ("service_id");
CREATE INDEX ON "Orders" ("status");
CREATE INDEX ON "Subscriptions" ("provider_id");
CREATE INDEX ON "Subscriptions" ("status");
CREATE INDEX ON "PromoCodes" ("partner_id");
CREATE INDEX ON "PromoCodes" ("code");
CREATE INDEX ON "Reviews" ("order_id");
CREATE INDEX ON "Reviews" ("client_id");
CREATE INDEX ON "Reviews" ("provider_id");
CREATE INDEX ON "Payment" ("user_id");
CREATE INDEX ON "Payment" ("status");
ALTER TABLE "Service"
ADD FOREIGN KEY ("provider_id") REFERENCES "Users" ("id");
ALTER TABLE "Service"
ADD FOREIGN KEY ("category_id") REFERENCES "ServiceCategories" ("id");
ALTER TABLE "Orders"
ADD FOREIGN KEY ("client_id") REFERENCES "Users" ("id");
ALTER TABLE "Orders"
ADD FOREIGN KEY ("service_id") REFERENCES "Service" ("id");
ALTER TABLE "Subscriptions"
ADD FOREIGN KEY ("provider_id") REFERENCES "Users" ("id");
ALTER TABLE "PromoCodes"
ADD FOREIGN KEY ("partner_id") REFERENCES "Users" ("id");
ALTER TABLE "Reviews"
ADD FOREIGN KEY ("order_id") REFERENCES "Orders" ("id");
ALTER TABLE "Reviews"
ADD FOREIGN KEY ("client_id") REFERENCES "Users" ("id");
ALTER TABLE "Reviews"
ADD FOREIGN KEY ("provider_id") REFERENCES "Users" ("id");
ALTER TABLE "Payment"
ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");