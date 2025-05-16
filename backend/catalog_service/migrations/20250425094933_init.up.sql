CREATE TABLE
    product (
        product_id UUID PRIMARY KEY NOT NULL,
        product_name VARCHAR(120) UNIQUE NOT NULL,
        description TEXT,
        price NUMERIC(10, 2) NOT NULL CHECK (price > 0),
        created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
        updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
    );

CREATE TABLE
    product_image (
        image_id UUID PRIMARY KEY NOT NULL,
        product_id UUID NOT NULL,
        url VARCHAR(255) NOT NULL,
        alt VARCHAR(255) NOT NULL,
        object_key TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
        updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
        CONSTRAINT fk_product_image FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
    );

CREATE TABLE
    category (
        category_id UUID PRIMARY KEY NOT NULL,
        category_name VARCHAR(100) UNIQUE NOT NULL,
        parent_id UUID DEFAULT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
        updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
        CONSTRAINT fk_category_parent FOREIGN KEY (parent_id) REFERENCES category (category_id) ON DELETE CASCADE
    );

CREATE TABLE
    product_category (
        product_id UUID NOT NULL,
        category_id UUID NOT NULL,
        assigned_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
        CONSTRAINT pk_product_category PRIMARY KEY (product_id, category_id),
        CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category (category_id) ON DELETE CASCADE,
        CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
    );