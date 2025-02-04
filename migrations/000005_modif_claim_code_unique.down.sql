ALTER TABLE claim_requests
DROP CONSTRAINT claim_code_unique;

ALTER TABLE claim_requests
ALTER COLUMN claim_code DROP NOT NULL;