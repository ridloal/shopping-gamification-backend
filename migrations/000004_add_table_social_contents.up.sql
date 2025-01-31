CREATE TYPE platform_enum AS ENUM ('facebook', 'instagram', 'twitter', 'youtube', 'tiktok');

CREATE TABLE social_contents (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL DEFAULT 0,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    platform platform_enum NOT NULL DEFAULT 'instagram',
    post_url VARCHAR(255) NOT NULL,
    status BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);