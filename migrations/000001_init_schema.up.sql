CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    image_url VARCHAR(255),
    stock INT NOT NULL DEFAULT 0,
    status BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE prizes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    discount_percentage INT NOT NULL,
    quota INT NOT NULL,
    remaining_quota INT NOT NULL,
    valid_until TIMESTAMPTZ,
    status BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE prize_groups (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    prize_id INT NOT NULL,
    probability NUMERIC(5,2) NOT NULL,
    status BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (prize_id) REFERENCES prizes(id)
);

CREATE TABLE claim_requests (
    id SERIAL PRIMARY KEY,
    user_id INT,
    product_id INT NOT NULL,
    prize_id INT,
    social_media_username VARCHAR(100) NOT NULL,
    social_media_platform VARCHAR(50) NOT NULL,
    post_url VARCHAR(255) NOT NULL,
    verification_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    claim_code VARCHAR(50) UNIQUE,
    claimed_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (prize_id) REFERENCES prizes(id)
);