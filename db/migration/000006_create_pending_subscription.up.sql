CREATE TABLE IF NOT EXISTS pending_subscriptions (
    id bigserial PRIMARY KEY,
    payment_id bigint NOT NULL REFERENCES payments(id),
    user_id bigint NOT NULL REFERENCES users(id),
    subscription_type varchar NOT NULL,
    start_date timestamptz NOT NULL,
    end_date timestamptz NOT NULL,
    original_price varchar NOT NULL,
    final_price varchar NOT NULL,
    promo_code_id bigint REFERENCES promo_codes(id),
    is_update boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);