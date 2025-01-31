ALTER TABLE products
DROP COLUMN original_price,
DROP COLUMN stars,
DROP COLUMN sold,
DROP COLUMN reviews,
DROP COLUMN external_link;
DROP COLUMN is_digital;

ALTER TABLE prize_groups
DROP COLUMN detail_json;

ALTER TABLE prizes
DROP COLUMN image_url;