ALTER TABLE api_keys ADD COLUMN key_prefix VARCHAR(20);

CREATE INDEX idx_api_keys_prefix ON api_keys(key_prefix) WHERE is_active = TRUE;
