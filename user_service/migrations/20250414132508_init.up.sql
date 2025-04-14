CREATE TABLE
    profile (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        avatar_url TEXT DEFAULT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL,
        updated_at TIMESTAMP DEFAULT NOW () NOT NULL
    );