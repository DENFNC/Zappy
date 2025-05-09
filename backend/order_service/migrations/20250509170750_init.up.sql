CREATE TABLE IF NOT EXISTS
    orders (
      order_id UUID PRIMARY KEY NOT NULL,
      profile_id UUID NOT NULL,
      status SMALLINT NOT NULL,
      total_amount_cents BIGINT NOT NULL,
      currency CHAR(3) NOT NULL,
      shipping_snapshot JSONB NOT NULL,
      payment_snapshot JSONB NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    )

CREATE TABLE IF NOT EXISTS
    order_items (
      order_item_id UUID PRIMARY KEY,
      order_id UUID NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
      product_id UUID NOT NULL,
      quantity INT NOT NULL,
      unit_price_cents BIGINT NOT NULL,
      currency CHAR(3) NOT NULL
    );


CREATE TABLE IF NOT EXISTS
    saga_state (
      saga_id UUID PRIMARY KEY,
      order_id UUID NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
      step SMALLINT NOT NULL,
      data JSONB,
      updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS
    outbox (
      id BIGSERIAL PRIMARY KEY,
      aggregate_id UUID NOT NULL,
      event_type TEXT NOT NULL,
      payload JSONB NOT NULL,
      published BOOLEAN NOT NULL DEFAULT FALSE,
      created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS
    inbox (
      id BIGSERIAL PRIMARY KEY,
      message_id UUID NOT NULL UNIQUE,
      event_type TEXT NOT NULL,
      payload JSONB NOT NULL,
      processed BOOLEAN NOT NULL DEFAULT FALSE,
      created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
