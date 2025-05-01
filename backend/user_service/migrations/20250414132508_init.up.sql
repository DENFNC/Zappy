-- 1. profile
CREATE TABLE
    profile (
        profile_id UUID PRIMARY KEY NOT NULL,
        auth_user_id UUID UNIQUE NOT NULL,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- 2. shipping_address
CREATE TABLE
    shipping_address (
        address_id UUID PRIMARY KEY NOT NULL,
        profile_id UUID NOT NULL REFERENCES profile (profile_id) ON DELETE CASCADE,
        country VARCHAR(50) NOT NULL,
        city VARCHAR(50) NOT NULL,
        street VARCHAR(100) NOT NULL,
        postal_code VARCHAR(20) NOT NULL,
        is_default BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Только один дефолтный адрес на профиль
CREATE UNIQUE INDEX ux_shipping_default ON shipping_address (profile_id)
WHERE
    is_default;

-- 3. payment_method
CREATE TABLE
    payment_method (
        payment_id UUID PRIMARY KEY NOT NULL,
        profile_id UUID NOT NULL REFERENCES profile (profile_id) ON DELETE CASCADE,
        payment_token VARCHAR(255) NOT NULL,
        is_default BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Только один дефолтный способ оплаты на профиль
CREATE UNIQUE INDEX ux_payment_default ON payment_method (profile_id)
WHERE
    is_default;

-- 4. wishlist_item
CREATE TABLE
    wishlist_item (
        item_id UUID PRIMARY KEY NOT NULL,
        profile_id UUID NOT NULL REFERENCES profile (profile_id) ON DELETE CASCADE,
        product_id UUID NOT NULL,
        added_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Запрет дублирования одного и того же товара для одного профиля
CREATE UNIQUE INDEX ux_wishlist_profile_product ON wishlist_item (profile_id, product_id);