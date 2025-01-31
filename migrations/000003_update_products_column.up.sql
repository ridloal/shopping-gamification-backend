ALTER TABLE products
ADD COLUMN original_price DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
ADD COLUMN stars DECIMAL(2, 1) DEFAULT 0.0,
ADD COLUMN sold INT DEFAULT 0,
ADD COLUMN reviews INT DEFAULT 0,
ADD COLUMN external_link VARCHAR(255),
ADD COLUMN is_digital BOOLEAN DEFAULT FALSE;

ALTER TABLE prize_groups
ADD COLUMN detail_json JSONB DEFAULT '{}'::jsonb;

ALTER TABLE prizes
ADD COLUMN image_url VARCHAR(255);