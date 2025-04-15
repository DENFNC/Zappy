-- 2.1. profiles
CREATE TABLE
    profiles (
        profile_id SERIAL PRIMARY KEY,
        auth_user_id INTEGER UNIQUE NOT NULL,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        phone VARCHAR(20),
        created_at TIMESTAMP DEFAULT NOW (),
        updated_at TIMESTAMP DEFAULT NOW ()
    );

-- 2.2. shipping_addresses
CREATE TABLE
    shipping_addresses (
        address_id SERIAL PRIMARY KEY,
        profile_id INTEGER NOT NULL REFERENCES profiles (profile_id),
        country VARCHAR(50) NOT NULL,
        city VARCHAR(50) NOT NULL,
        street VARCHAR(100) NOT NULL,
        postal_code VARCHAR(20) NOT NULL,
        is_default BOOLEAN DEFAULT FALSE
    );

-- 2.3. payment_methods
CREATE TABLE
    payment_methods (
        payment_id SERIAL PRIMARY KEY,
        profile_id INTEGER NOT NULL REFERENCES profiles (profile_id),
        payment_token VARCHAR(255) NOT NULL,
        is_default BOOLEAN DEFAULT FALSE
    );

-- 2.4. wishlist_items (привязка к profile_id напрямую)
CREATE TABLE
    wishlist_items (
        item_id SERIAL PRIMARY KEY,
        profile_id INTEGER NOT NULL REFERENCES profiles (profile_id),
        product_id INTEGER NOT NULL,
        added_at TIMESTAMP DEFAULT NOW ()
    );