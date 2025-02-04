ALTER TABLE claim_requests
ADD COLUMN prize_detail JSONB DEFAULT '{}'::jsonb;