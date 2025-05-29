CREATE TABLE
    review (
        review_id UUID PRIMARY KEY,
        product_id UUID NOT NULL,
        profile_id UUID NOT NULL,
        rating SMALLINT NOT NULL,
        COMMENTS TEXT,
        created_at TIMESTAMP DEFAULT NOW() NOT NULL,
        updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
        CONSTRAINT CHECK rating_check (rating BETWEEN 1 AND 5)
    );