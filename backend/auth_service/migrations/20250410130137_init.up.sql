CREATE TABLE
    users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        username VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL,
        updated_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    roles (
        id SERIAL PRIMARY KEY,
        role_name VARCHAR(20),
        created_at TIMESTAMP DEFAULT NOW () NOT NULL,
        updated_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    user_roles (
        user_id BIGINT,
        role_id BIGINT,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL,
        CONSTRAINT pk_user_role PRIMARY KEY (user_id, role_id),
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
        CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
    );