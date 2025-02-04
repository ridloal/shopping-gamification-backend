ALTER TABLE claim_requests
ALTER COLUMN claim_code SET NOT NULL;
ALTER TABLE claim_requests
ADD CONSTRAINT claim_code_unique UNIQUE (claim_code);
